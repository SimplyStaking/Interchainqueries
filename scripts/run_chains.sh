echo "Starting Source Chain in new terminal"

gnome-terminal -- starport chain serve -c config-icq-1.yml --reset-once -v

echo "Starting Target Chain in new terminal"

sleep 30s

gnome-terminal -- starport chain serve -c config-icq-2.yml --reset-once