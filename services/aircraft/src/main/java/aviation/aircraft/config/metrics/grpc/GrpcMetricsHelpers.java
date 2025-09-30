package aviation.aircraft.config.metrics.grpc;

import io.grpc.MethodDescriptor;
import io.grpc.Status;
import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.Meter;
import io.micrometer.core.instrument.Timer;

/**
 * A helper class for recording metrics for gRPC requests.
 */
public class GrpcMetricsHelpers {

  /**
   * A record for a gRPC method.
   *
   * @param service the service the method belongs to.
   * @param method  the method name.
   */
  public record MethodParts(String service, String method) {
  }

  /**
   * Extract the service and method from a full method name.
   *
   * @param fullMethod the full method name.
   *
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
   * @param counterProvider the counter-provider.
   * @param timerProvider   the timer provider.
   * @param service         the service the method belongs to.
   * @param method          the method name.
   * @param code            the status code of the request.
   * @param sample          the sample to stop the timer with.
   */
  public static void recordMetrics(
          Meter.MeterProvider<Counter> counterProvider,
          Meter.MeterProvider<Timer> timerProvider,
          String service,
          String method,
          Status.Code code,
          Timer.Sample sample) {
    Counter counter =
            counterProvider.withTags("service", service,
                    "method", method,
                    "code", code.name());
    Timer timer =
            timerProvider.withTags("service", service,
                    "method", method,
                    "code", code.name());
    counter.increment();
    sample.stop(timer);
  }
}