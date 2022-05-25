echo "Starting Source Chain in new terminal"

gnome-terminal -- starport chain serve -c config-ica-1.yml --reset-once

echo "Starting Target Chain in new terminal"

sleep 30s

gnome-terminal -- starport chain serve -c config-ica-2.yml --reset-once