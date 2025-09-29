package aviation.aircraft.config.metrics.grpc;

import io.grpc.MethodDescriptor;
import io.grpc.Status;
import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;

/**
 * A helper class for recording metrics for gRPC requests.
 */
public class GrpcMetricsHelpers {

  /**
   * A record for a gRPC method.
   *
   * @param service the service the method belongs to.
   * @param method the method name.
   */
  public record MethodParts(String service, String method) {}

  /**
   * Extract the service and method from a full method name.
   *
   * @param fullMethod the full method name.
   * @return the service and method.
   */
  public static MethodParts extract(String fullMethod) {
    String service = MethodDescriptor.extractFullServiceName(fullMethod);
    if (service == null) {
      service = "unknown";
    }
    String method = fullMethod.substring(fullMethod.lastIndexOf('/') + 1);
    return new MethodParts(service, method);
  }

  /**
   * Record metrics for a gRPC request.
   *
   * @param registry the registry to record metrics in.
   * @param counterBuilder the counter-builder to use.
   * @param timerBuilder the timer builder to use.
   * @param service the service the method belongs to.
   * @param method the method name.
   * @param code the status code of the request.
   * @param sample the sample to stop the timer with.
   */
  public static void recordMetrics(
          MeterRegistry registry,
          Counter.Builder counterBuilder,
          Timer.Builder timerBuilder,
          String service,
          String method,
          Status.Code code,
          Timer.Sample sample) {

    Counter counter = counterBuilder
            .tags("service", service, "method", method, "code", code.name())
            .register(registry);

    Timer timer = timerBuilder
            .tags("service", service, "method", method, "code", code.name())
            .register(registry);

    counter.increment();
    sample.stop(timer);
  }
}