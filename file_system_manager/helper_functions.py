import os

def get_image_files(path, image_data_func):
    # List of common image file extensions, including RAW formats
    image_extensions = {
        '.jpg', '.jpeg', 
        '.cr2', '.cr3', '.raw'
    }
    
    # Check if the provided path exists and is a directory
    if not os.path.isdir(path):
        raise ValueError(f"The path '{path}' is not a valid directory.")
    
    # Initialize the result dictionary
    image_files_data = {}

    # Get all files in the directory and filter by image extension
    for file in os.listdir(path):
        file_path = os.path.join(path, file)
        
        if os.path.isfile(file_path) and any(file.lower().endswith(ext) for ext in image_extensions):
            # Call the image_data_func for the image and store the result
            image_data = image_data_func(file_path)
            
            # Store the image data in the dictionary
            image_files_data[file] = image_data
    
    return image_files_data

def get_image_files_in_subfolders(path, image_data_func):
    # Check if the provided path exists and is a directory
    if not os.path.isdir(path):
        raise ValueError(f"The path '{path}' is not a valid directory.")
    
    all_images = {}  # Dictionary to store images found in each subfolder
    
    # Walk through the directory and its subdirectories
    for root, dirs, files in os.walk(path):
        # Get image files for the current folder
        image_files = get_image_files(root, image_data_func)
        
        if image_files:
            all_images[root] = image_files
    
    return all_images

