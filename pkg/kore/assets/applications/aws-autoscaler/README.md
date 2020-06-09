# Update manifests

**TODO** Update build to:

1. Add kustomize? to edit files downloaded,
1. Add a download per k8 version in line with supported k8 versions

For now, download version:

1. Download the manifests from the kubernetes/autoscaler Github repo:

    ```
    (
    cd pkg/kore/assets/applications/aws-autoscaler
    K8S_AUTOSCALER_VERSION=1.16.5
    curl -sSL "https://raw.githubusercontent.com/kubernetes/autoscaler/cluster-autoscaler-${K8S_AUTOSCALER_VERSION}/cluster-autoscaler/cloudprovider/aws/examples/cluster-autoscaler-multi-asg.yaml" -o cluster-autoscaler-multi-asg.yaml
    )
    ```

1. Edit the file to add the following changes:

    ```
    kind: ServiceAccount
    metadata:
      annotaions:
        eks.amazonaws.com/role-arn: {{ .AutoscalingRoleARN }}
    ```

    ```
    kind: Deployment
    ...
       - image: {{ .AutoScalingImage }}
    ...
         command:
    ...
            {{- range $key, $value := .ASGs }}
            - --nodes={{ $.value.Min }}:{{ $.value.Max }}:{{ $.value.Name }}
            {{- end }}
          env:
            - name: AWS_REGION
              value: {{ .AwsRegion }}
    ```
