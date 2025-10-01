package aviation.aircraft.config.metrics.grpc;

import io.grpc.ForwardingServerCall;
import io.grpc.Metadata;
import io.grpc.ServerCall;
import io.grpc.ServerCallHandler;
import io.grpc.ServerInterceptor;
import io.grpc.Status;
import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.common.Attributes;
import io.opentelemetry.api.metrics.DoubleHistogram;
import io.opentelemetry.api.metrics.LongCounter;
import io.opentelemetry.api.metrics.Meter;
import net.devh.boot.grpc.server.interceptor.GrpcGlobalServerInterceptor;
import org.springframework.stereotype.Component;

/**
 * A gRPC server interceptor that records metrics for inbound gRPC requests.
 */
@Component
@GrpcGlobalServerInterceptor
public class GrpcInboundServerMetricsInterceptor implements ServerInterceptor {

  private static final String SERVICE = "service";
  private static final String METHOD = "method";
  private static final String STATUS = "status";

  private final LongCounter requestCounter;
  private final DoubleHistogram requestDuration;

  /**
   * A constructor for the interceptor.
   */
  public GrpcInboundServerMetricsInterceptor() {
    Meter meter = GlobalOpenTelemetry.getMeter("aviation.aircraft.grpc");

    this.requestCounter = meter.counterBuilder("aircraft_grpc_requests_total")
            .setDescription("Total inbound gRPC requests for the aircraft service")
            .setUnit("1")
            .build();

    this.requestDuration = meter.histogramBuilder("aircraft_grpc_request_duration_ms")
            .setDescription("Duration of inbound gRPC requests")
            .setUnit("ms")
            .build();
  }

  /**
   * Intercepts the call and records metrics.
   *
   * @param call    the call to intercept.
   * @param headers the headers of the call.
   * @param next    the next handler in the chain.
   * @param <ReqT>  the request type.
   * @param <RespT> the response type.
   *
   * @return the listener for the call.
   */
  @Override
  public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(
          ServerCall<ReqT, RespT> call,
          Metadata headers,
          ServerCallHandler<ReqT, RespT> next) {

    long start = System.nanoTime();
    GrpcMetricsHelpers.MethodParts parts =
            GrpcMetricsHelpers.extract(call.getMethodDescriptor().getFullMethodName());

    ServerCall<ReqT, RespT> monitoringCall =
            new ForwardingServerCall.SimpleForwardingServerCall<>(call) {
              @Override
              public void close(Status status, Metadata trailers) {
                long durationMs = (System.nanoTime() - start) / 1_000_000;

                Attributes attrs = Attributes.builder()
                        .put(SERVICE, parts.service())
                        .put(METHOD, parts.method())
                        .put(STATUS, status.getCode().name())
                        .build();

                requestCounter.add(1, attrs);
                requestDuration.record(durationMs, attrs);

                super.close(status, trailers);
              }
            };

    return next.startCall(monitoringCall, headers);
  }
}