echo "Attempting to delete path, ignore if errors"

rly paths delete stpath

echo "Creating a path and linking"

rly paths new source-chain target-chain stpath
rly transact link stpath
rly paths icq stpath true 15

echo "Starting relayer in new terminal"

gnome-terminal -- rly start stpath