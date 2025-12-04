from dotenv import load_dotenv
import os


def load_config(config_file: str) -> dict[str, str|bool]:
    if not config_file:
        load_dotenv()
    else:     
        load_dotenv(dotenv_path=config_file)
    config: dict[str, str|bool] = {}
    config["host"] = os.getenv("PMG_HOST", "default_host")
    config["user"] = os.getenv("PMG_USER", "default_user")
    config["password"] = os.getenv("PMG_PASSWORD", "default_password")
    config["verify_ssl"] = os.getenv("PMG_VERIFY_SSL", "true").lower() == "true"
    config["backend"] = os.getenv("PMG_BACKEND", "https")
    config["service"] = os.getenv("PMG_SERVICE", "pmg")
    config["exporter_port"] = os.getenv("PMG_EXPORTER_PORT", "10069")
    config["log_level"] = os.getenv("PMG_LOG_LEVEL", "INFO")
    return config
# 