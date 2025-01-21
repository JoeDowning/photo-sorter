from PIL import Image
from PIL.ExifTags import TAGS
from rawkit.raw import Raw
import os

def extract_exif_data(image_path):
    _, ext = os.path.splitext(image_path)
    ext = ext.lower()

    if ext == '.cr3' or ext == '.cr2' or ext == '.raw':
        with Raw(filename=image_path) as raw:
            # Extract EXIF data from raw image
            metadata = raw.metadata
            date_taken = metadata.timestamp
            camera_model = metadata.camera

            return [date_taken, camera_model]
    else:
        # Open the image using Pillow
        with Image.open(image_path) as img:
            # Try to get the EXIF data
            exif_data = img._getexif()
            
            if exif_data is not None:
                # Extract the DateTimeOriginal and Model from EXIF data
                date_taken = None
                camera_model = None
                
                for tag, value in exif_data.items():
                    tag_name = TAGS.get(tag, tag)
                    
                    if tag_name == 'DateTimeOriginal':
                        date_taken = value
                    elif tag_name == 'Model':
                        camera_model = value
                
                return [date_taken, camera_model]
            else:
                # If no EXIF data is found, return None for both
                return [None, None]