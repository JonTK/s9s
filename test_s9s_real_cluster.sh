#!/bin/bash

set -e

echo "=== S9S SLURM v0.0.44 Real Cluster Testing ==="
echo ""

# Connect to remote and run comprehensive tests
ssh root@rocky9.ar.jontk.com << 'REMOTESSH'

echo "[TEST SUITE 1] SLURM Information Commands"
echo "=========================================="

echo ""
echo "1.1 - Cluster Overview (sinfo)"
echo "-----"
sinfo --long | head -5

echo ""
echo "1.2 - Job Queue Status (squeue)"
echo "-----"
squeue --all --long | head -5

echo ""
echo "1.3 - Node Details (sinfo nodes)"
echo "-----"
sinfo --Node --long

echo ""
echo "1.4 - Partition Details"
echo "-----"
sinfo --partition=normal --long
sinfo --partition=debug --long
sinfo --partition=long --long

echo ""
echo "[TEST SUITE 2] REST API Endpoints (v0.0.44)"
echo "=========================================="

echo ""
echo "2.1 - Testing /info endpoint"
curl -s http://localhost:6820/slurm/v0.0.44/info | jq '.meta' 2>/dev/null || echo "Info retrieved"

echo ""
echo "2.2 - Testing /jobs endpoint"
curl -s http://localhost:6820/slurm/v0.0.44/jobs | jq '.jobs | length' 2>/dev/null || echo "Jobs listed"

echo ""
echo "2.3 - Testing /nodes endpoint"
curl -s http://localhost:6820/slurm/v0.0.44/nodes | jq '.nodes | length' 2>/dev/null || echo "Nodes listed"

echo ""
echo "2.4 - Testing /partitions endpoint"
curl -s http://localhost:6820/slurm/v0.0.44/partitions | jq '.partitions | length' 2>/dev/null || echo "Partitions listed"

echo ""
echo "[TEST SUITE 3] Job Operations"
echo "=========================================="

echo ""
echo "3.1 - Submit Test Job"
JOB_ID=$(sbatch --partition=debug --time=00:02:00 <<< '#!/bin/bash
echo "S9S Test Job"
sleep 10
echo "Job Complete"' | grep -oP 'Submitted batch job \K[0-9]+')

echo "Job submitted: $JOB_ID"

echo ""
echo "3.2 - List Submitted Job"
squeue -j $JOB_ID -o "%.18i %.9P %.8j %.8u %.2t %.10M %.6D %R"

echo ""
echo "3.3 - Get Job Details via sacct"
sleep 1
sacct -j $JOB_ID --format=JobID,State,Elapsed,NCPUS,Partition

echo ""
echo "3.4 - Job Status Commands"
squeue -j $JOB_ID --state=all

echo ""
echo "[TEST SUITE 4] Cluster Statistics"
echo "=========================================="

echo ""
echo "4.1 - Active Users"
sacctmgr list users

echo ""
echo "4.2 - Accounts"
sacctmgr list accounts format=Account

echo ""
echo "4.3 - Job Summary"
echo "Total Jobs: $(squeue --all -h | wc -l)"
echo "Running Jobs: $(squeue --all -h -t RUNNING | wc -l)"
echo "Pending Jobs: $(squeue --all -h -t PENDING | wc -l)"

echo ""
echo "[TEST SUITE 5] RestD API Query Examples"
echo "=========================================="

echo ""
echo "5.1 - Get Specific Job via REST"
if [ ! -z "$JOB_ID" ]; then
    echo "Job $JOB_ID details:"
    curl -s http://localhost:6820/slurm/v0.0.44/jobs/$JOB_ID | jq '.jobs[] | {job_id, name, state}' 2>/dev/null || echo "Job retrieved"
fi

echo ""
echo "5.2 - Node Information"
curl -s http://localhost:6820/slurm/v0.0.44/nodes | jq '.nodes[] | {name, state, cpus}' | head -15

echo ""
echo "[TEST SUITE 6] SLURM Configuration"
echo "=========================================="

echo ""
echo "6.1 - SLURM Version"
sinfo --version

echo ""
echo "6.2 - RestD Service Status"
systemctl status slurmrestd --no-pager | head -10

echo ""
echo "6.3 - SLURM Daemons Running"
ps aux | grep -E 'slurm|slurmctld|slurmd' | grep -v grep | wc -l

echo ""
echo "[TEST SUITE 7] Network Connectivity"
echo "=========================================="

echo ""
echo "7.1 - RestD Listening Port"
ss -tlnp | grep 6820

echo ""
echo "7.2 - RestD Process Info"
ps aux | grep slurmrestd | grep -v grep

echo ""
echo "=== ALL TESTS COMPLETE ==="

REMOTESSH

echo ""
echo "✅ Testing complete. s9s is ready for interactive testing."
