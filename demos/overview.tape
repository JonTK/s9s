Require s9s

Output demos/output/overview.gif
Output demos/output/overview.mp4

Source demos/common.tape

# Setup environment
Hide
Type "export S9S_ENABLE_MOCK=dev"
Enter
Sleep 500ms
Type "clear"
Enter
Sleep 500ms
Show

# Launch s9s
Type "./build/s9s"
Sleep 500ms
Enter
Sleep 2s

# Show Jobs view (default - key 1)
Sleep 2s

# Navigate to Nodes (key 2)
Type "2"
Sleep 1500ms

# Navigate to Partitions (key 3)
Type "3"
Sleep 1500ms

# Navigate to Reservations (key 4)
Type "4"
Sleep 1500ms

# Navigate to QoS (key 5)
Type "5"
Sleep 1500ms

# Navigate to Accounts (key 6)
Type "6"
Sleep 1500ms

# Navigate to Users (key 7)
Type "7"
Sleep 1500ms

# Navigate to Dashboard (key 8)
Type "8"
Sleep 2s

# Navigate to Health (key 9)
Type "9"
Sleep 2s

# Show Global Search (Ctrl+F)
Ctrl+F
Sleep 1s
Type "compute"
Sleep 1500ms
Escape
Sleep 500ms

# Back to Jobs (key 1)
Type "1"
Sleep 1s

# Show filtering (simple text filter)
Type "/"
Sleep 500ms
Type "RUNNING"
Sleep 1s
Enter
Sleep 1500ms

# Clear filter and exit
Escape
Sleep 500ms
Type "q"
Sleep 500ms
