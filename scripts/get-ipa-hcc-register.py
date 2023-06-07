#!/usr/bin/python3
"""Ephemeral: Create new stub domain and print ipa-hcc register command
"""
import base64
import json
import subprocess

import requests
import requests.auth

CREATE_JSON = {
    "title": "Human readable title",
    "description": "My human friendly description",
    "auto_enrollment_enabled": True,
    "domain_type": "rhel-idm",
}


def oc(*args) -> str:
    cmd = ["oc"]
    cmd.extend(args)
    return subprocess.check_output(cmd, text=True).strip()


def main() -> None:
    namespace = oc("project", "-q")

    keycloak = oc("get", f"secrets/env-{namespace}-keycloak", "-o", "json")
    secrets = json.loads(keycloak)
    username = base64.b64decode(secrets["data"]["defaultUsername"]).decode("utf-8")
    password = base64.b64decode(secrets["data"]["defaultPassword"]).decode("utf-8")

    hmsidm_backend = oc(
        "get",
        "routes",
        "-l",
        "app=hmsidm-backend",
        "-o",
        "jsonpath={.items[0].spec.host}",
    )

    url = f"https://{hmsidm_backend}/api/hmsidm/v1/domains"
    print(f"POST {url}")
    resp = requests.post(
        url,
        auth=requests.auth.HTTPBasicAuth(username, password),
        headers={
            "X-Rh-Insights-Request-Id": "test_12345",
        },
        json=CREATE_JSON,
    )

    domain_id = resp.json()["domain_id"]
    hdr = resp.headers["x-rh-idm-rhelidm-register-token"]
    token = json.loads(base64.b64decode(hdr))
    domain_secret = token["secret"]

    print(f"ipa-hcc register {domain_id} {domain_secret}")


if __name__ == "__main__":
    main()
