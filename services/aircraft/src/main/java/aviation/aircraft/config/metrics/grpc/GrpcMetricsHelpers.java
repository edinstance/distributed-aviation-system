package aviation.aircraft.config.metrics.grpc;

import io.grpc.MethodDescriptor;

/**
 * A helper class for the grpc metrics.
 */
public class GrpcMetricsHelpers {

  /**
   * A record for the service and method.
   *
   * @param service the service.
   * @param method the method.
   */
  public record MethodParts(String service, String method) {}

  /**
   * Extracts the service and method from a full method name.
   *
   * @param fullMethod the full method name.
   * @return the service and method.
   */
  public static MethodParts extract(String fullMethod) {
    String service = MethodDescriptor.extractFullServiceName(fullMethod);
    if (service == null || service.isBlank()) {
      service = "unknown";
    }
    String method = fullMethod.substring(fullMethod.lastIndexOf('/') + 1);
    if (method.isBlank()) {
      method = "unknown";
    }
    return new MethodParts(service, method);
  }
}