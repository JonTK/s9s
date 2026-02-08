# Bootstrap script for s9s demos
# Source this after common.tape to set up environment and launch s9s

# Setup environment (hidden from recording)
Hide
Type "export S9S_ENABLE_MOCK=dev"
Enter
Sleep 500ms
Type "clear"
Enter
Sleep 500ms
Show

# Launch s9s
Type "./s9s"
Enter
Sleep 2s
