echo "Uninstalling Peek!"
sudo rm /usr/local/bin/peek # Remove the executable

# ask user if they want to remove the config file
# shellcheck disable=SC2162
read -p "Do you want to remove the config file? [y/n]: " remove_config

if [ "$remove_config" == "y" ]; then
    echo "Removing the config file.."
    # get home directory
    home_dir=$(eval echo ~"$USER")
    # delete peek config file
    echo "Deleting \"$home_dir/.config/peek\""
    sudo rm -rf "$home_dir"/.config/peek
fi
