# News Aggregator Operator

The News Aggregator Operator manages the lifecycle of news feed resources in a Kubernetes cluster. It automates the addition, update, and deletion of news sources by interacting with a backend service that handles news aggregation.

## Description

This operator is designed to manage `Feed` resources, which represent individual news sources. It provides functionality to add, update, and remove these sources, ensuring that the backend service remains synchronized with the state of the `Feed` resources in your Kubernetes cluster.

## Getting Started

### Prerequisites

Before deploying the operator, ensure you have the following tools installed:

- **Go**: Version 1.22.0+
- **Docker**: Version 17.03+
- **kubectl**: Version 1.11.3+
- **Kubernetes Cluster**: Version 1.11.3+ (Ensure you have access to a Kubernetes cluster)

### Deployment

1. **Build and Push the Docker Image**

   Build and push your operator image to your container registry:

   ```sh
   make docker-build docker-push IMG=<some-registry>/operator:tag
   ```

   **Note:** Ensure the image is published in the registry specified and that you have access to pull the image.

2. **Install Custom Resource Definitions (CRDs)**

   Apply the Custom Resource Definitions (CRDs) to your cluster:

   ```sh
   make install
   ```

3. **Deploy the Operator**

   Deploy the operator using the image built in the previous step:

   ```sh
   make deploy IMG=<some-registry>/operator:tag
   ```

   **Note:** If you encounter RBAC (Role-Based Access Control) errors, you might need to grant cluster-admin privileges or ensure you are logged in as an admin.

4. **Create `Feed` Instances**

   Apply sample `Feed` resources to test the operator:

   ```sh
   kubectl apply -k config/samples/
   ```

   **Note:** Verify that the samples have default values set for proper testing.

### Uninstalling

1. **Delete Instances**

   Remove the `Feed` resources from the cluster:

   ```sh
   kubectl delete -k config/samples/
   ```

2. **Remove CRDs**

   Uninstall the CRDs:

   ```sh
   make uninstall
   ```

3. **Undeploy the Operator**

   Remove the operator from the cluster:

   ```sh
   make undeploy
   ```

## Project Distribution

To distribute the operator, follow these steps:

1. **Build the Installer**

   Build the installer for your operator image:

   ```sh
   make build-installer IMG=<some-registry>/operator:tag
   ```

   **Note:** This generates an `install.yaml` file in the `dist` directory, which contains the necessary resources to install the operator.

---

### CRD Definitions

#### `Feed` Custom Resource

**Spec:**
- `name` (string): The name of the news source.
- `link` (string): The URL of the news feed.

**Status:**
- `conditions` (array of `Condition`): The status conditions of the `Feed`.

**Condition Types:**
- `Added`: The feed has been added successfully.
- `Updated`: The feed has been updated successfully.
- `Deleted`: The feed has been deleted successfully.

**Condition Structure:**
- `type` (string): The type of condition.
- `status` (boolean): The status of the condition.
- `reason` (string, optional): The reason for the condition.
- `message` (string, optional): A message describing the condition.
- `lastUpdateTime` (Time): The last update time of the condition.

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
