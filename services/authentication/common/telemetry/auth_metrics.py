from common.telemetry.helpers import get_meter


class AuthMetrics:
    """Custom metrics for authentication service."""

    def __init__(self):
        self.meter = get_meter()

        # ---- Counters ----
        self.login_attempts = self.meter.create_counter(
            "auth_login_attempts_total",
            description="Total number of login attempts.",
            unit="1",
        )

        self.login_successes = self.meter.create_counter(
            "auth_login_successes_total",
            description="Total number of successful logins.",
            unit="1",
        )

        self.login_failures = self.meter.create_counter(
            "auth_login_failures_total",
            description="Total number of failed login attempts.",
            unit="1",
        )

        self.token_generations = self.meter.create_counter(
            "auth_token_generations_total",
            description="Total number of JWT tokens generated.",
            unit="1",
        )

        self.token_validations = self.meter.create_counter(
            "auth_token_validations_total",
            description="Total number of JWT token validations.",
            unit="1",
        )

        # ---- Histograms ----
        self.login_duration = self.meter.create_histogram(
            "auth_login_duration_seconds",
            description="Duration of login requests.",
            unit="s",
        )

        self.token_validation_duration = self.meter.create_histogram(
            "auth_token_validation_duration_seconds",
            description="Duration of token validation.",
            unit="s",
        )

    def record_login_attempt(self, tenant: str, success: bool, duration: float):
        """Record a login attempt."""
        labels = {"tenant": tenant, "success": str(success).lower()}
        self.login_attempts.add(1, labels)
        if success:
            self.login_successes.add(1, labels)
        else:
            self.login_failures.add(1, labels)
        self.login_duration.record(duration, labels)

    def record_token_generation(self, tenant: str, token_type: str):
        """Record token generation."""
        labels = {"tenant": tenant, "token_type": token_type}
        self.token_generations.add(1, labels)

    def record_token_validation(self, tenant: str, valid: bool, duration: float):
        """Record token validation."""
        labels = {"tenant": tenant, "valid": str(valid).lower()}
        self.token_validations.add(1, labels)
        self.token_validation_duration.record(duration, labels)