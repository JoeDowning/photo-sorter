import sys

# Add the root of your project (not the 'main' directory)
sys.path.append('/Users/joe.downing/go/src/github.com/photos-sorter')

from file_system_manager.helper_functions import get_image_files_in_subfolders
from image_data_importer.image_data_extractor import extract_exif_data

sys.path.append('/Users/joe.downing/go/src/github.com/photos-sorter/main')

# Example usage
path = '/Users/joe.downing/Pictures/photos/testing-folder/'
image_files_in_subfolders = get_image_files_in_subfolders(path, extract_exif_data)

# Print the result
for folder, images in image_files_in_subfolders.items():
    print(f"Images in folder '{folder}': {images}")