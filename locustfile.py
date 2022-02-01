from locust import HttpUser, task

class VCICheckUser(HttpUser):
    @task
    def vci_check(self):
        self.client.get("/?iss=asdf")
