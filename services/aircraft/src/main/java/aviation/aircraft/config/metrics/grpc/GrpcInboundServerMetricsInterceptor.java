package aviation.aircraft.config.metrics.grpc;

import io.grpc.ForwardingServerCall;
import io.grpc.ForwardingServerCallListener;
import io.grpc.Metadata;
import io.grpc.ServerCall;
import io.grpc.ServerCallHandler;
import io.grpc.ServerInterceptor;
import io.grpc.Status;
import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import net.devh.boot.grpc.server.interceptor.GrpcGlobalServerInterceptor;
import org.springframework.stereotype.Component;

/**
 * A gRPC server interceptor that records metrics for inbound gRPC requests.
 */
@Component
@GrpcGlobalServerInterceptor
public class GrpcInboundServerMetricsInterceptor implements ServerInterceptor {

  private final MeterRegistry registry;
  private final Counter.Builder requestCounter;
  private final Timer.Builder requestTimer;

  /**
   * Constructor for the interceptor.
   *
   * @param registry the registry to record metrics in.
   */
  public GrpcInboundServerMetricsInterceptor(MeterRegistry registry) {
    this.registry = registry;
    this.requestCounter = Counter.builder("aircraft_grpc_requests_total")
            .description("Total number of inbound gRPC requests for the aircraft service")
            .tag("direction", "inbound");
    this.requestTimer = Timer.builder("aircraft_grpc_request_duration_seconds")
            .description("Duration of inbound gRPC requests for the aircraft service")
            .tag("direction", "inbound");
  }

  /**
   * Intercepts the call and records metrics.
   *
   * @param call the call to intercept.
   * @param headers the headers of the call.
   * @param next the next handler in the chain.
   * @param <ReqT> the request type.
   * @param <RespT> the response type.
   * @return the listener for the call.
   */
  @Override
  public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(
          ServerCall<ReqT, RespT> call,
          Metadata headers,
          ServerCallHandler<ReqT, RespT> next) {

    final GrpcMetricsHelpers.MethodParts parts =
            GrpcMetricsHelpers.extract(call.getMethodDescriptor().getFullMethodName());
    final Timer.Sample sample = Timer.start(registry);

    ServerCall<ReqT, RespT> monitoringCall =
            new ForwardingServerCall.SimpleForwardingServerCall<>(call) {
              @Override
              public void close(Status status, Metadata trailers) {
                GrpcMetricsHelpers.recordMetrics(
                        registry,
                        requestCounter,
                        requestTimer,
                        parts.service(),
                        parts.method(),
                        status.getCode(),
                        sample);
                super.close(status, trailers);
              }
            };

    return new ForwardingServerCallListener.SimpleForwardingServerCallListener<>(
            next.startCall(monitoringCall, headers)) {};
  }
}