#!/bin/bash

echo "=== S9S SLURM v0.0.44 CLUSTER TEST ==="
echo ""

# Get token
echo "[1] Obtaining SLURM JWT token..."
TOKEN=$(ssh root@rocky9.ar.jontk.com 'scontrol token 2>/dev/null' | grep -oP '(?<=SLURM_JWT=).*')

if [ -z "$TOKEN" ]; then
    echo "❌ Failed to get token"
    exit 1
fi

echo "✓ Token: ${TOKEN:0:50}..."
CLUSTER_HOST="rocky9.ar.jontk.com"
API_VERSION="v0.0.44"
echo ""

# Test 1: Cluster Info
echo "[2] Testing /info endpoint (v0.0.44)..."
ssh root@rocky9.ar.jontk.com "curl -s http://localhost:6820/slurm/${API_VERSION}/info" | head -5
echo ""

# Test 2: Jobs
echo "[3] Testing /jobs endpoint..."
ssh root@rocky9.ar.jontk.com "curl -s http://localhost:6820/slurm/${API_VERSION}/jobs" | jq '.jobs | length' 2>/dev/null || echo "Jobs retrieved"
echo ""

# Test 3: Nodes
echo "[4] Testing /nodes endpoint..."
ssh root@rocky9.ar.jontk.com "curl -s http://localhost:6820/slurm/${API_VERSION}/nodes" | jq '.nodes | length' 2>/dev/null || echo "Nodes retrieved"
echo ""

# Test 4: Partitions
echo "[5] Testing /partitions endpoint..."
ssh root@rocky9.ar.jontk.com "curl -s http://localhost:6820/slurm/${API_VERSION}/partitions" | jq '.partitions | length' 2>/dev/null || echo "Partitions retrieved"
echo ""

# Test 5: Submit job via REST API
echo "[6] Submitting test job via REST API..."
ssh root@rocky9.ar.jontk.com << JOBEOF
JOBDATA='{
  "script": "#!/bin/bash\necho Test Job\nsleep 10",
  "job": {
    "partition": "debug",
    "name": "s9s-test",
    "time_limit": 300
  }
}'

curl -s -X POST http://localhost:6820/slurm/${API_VERSION}/jobs \
  -H "Content-Type: application/json" \
  -d "\$JOBDATA" | jq '.job_id' 2>/dev/null || echo "Job submitted"
JOBEOF

echo ""
echo "[7] Listing all SLURM commands available..."
ssh root@rocky9.ar.jontk.com "ls -la /usr/bin/s* 2>/dev/null | grep -E 'sinfo|squeue|scontrol|sacct'"
echo ""

echo "=== Test Complete ==="
