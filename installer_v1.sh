echo "Installing Peek!.."

echo "Cloning the repo.."
git clone https://github.com/fwuffyboi/peek.git peek || echo "Failed to clone repo." && exit
cd peek/src || echo "Failed to cd into directory."

echo "Building the application.."
go build -o peek . || echo "Failed to build application." && exit # Build the file

echo "Setting permissions for the executable.."
sudo chmod +x peek || echo "Could not change Peek permissions." && exit # Make the file executable

echo "Moving the file to /usr/local/bin.."
sudo mv peek /usr/local/bin/peek || echo "Failed to move Peek to /usr/local/bin/peek" && exit # Move the file to /usr/local/bin

echo "Cleaning up.."
cd ../.. # Get out of the directory
sudo rm -rf peek || echo "Failed to rm -rf peek repo." && exit # Delete the unnecessary repo

echo "Done! Peek is now installed. Run 'peek' to start the application."