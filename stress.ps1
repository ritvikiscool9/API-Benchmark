param (
    [Parameter(Mandatory=$true)] [string]$url,
    [Parameter(Mandatory=$true)] [int]$count
)

Write-Host "--- Starting Stress Test ---" -ForegroundColor Cyan
Write-Host "Target: $url"
Write-Host "Requests: $count"

# 1. Clean up old runs
Write-Host "Cleaning up previous jobs..."
kubectl delete job worker-job --ignore-not-found=true

# 2. Read the YAML, swap placeholders, and apply to K8s
(Get-Content worker.yaml) `
    -replace "TARGET_URL_PLACEHOLDER", $url `
    -replace "REQUEST_COUNT_PLACEHOLDER", $count | `
    kubectl apply -f -

Write-Host "--- Job Dispatched ---" -ForegroundColor Green
Write-Host "Run 'kubectl logs job/worker-job' to see progress."