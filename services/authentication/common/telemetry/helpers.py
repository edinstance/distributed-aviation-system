from opentelemetry import trace, metrics

def get_tracer(name: str = "authentication-service"):
    """Return a tracer instance for a specific component."""
    return trace.get_tracer(name)


def get_meter(name: str = "authentication-service"):
    """Return a meter instance for a specific component."""
    return metrics.get_meter(name)