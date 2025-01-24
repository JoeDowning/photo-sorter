# photo-sorter

A small project to help with sorting photos for backup purposes

 - imageFileTypes: list of file types to consider as images
 - sourcePath: path to the folder containing the images
 - destinationPath: path to the folder where the images will be copied to

The script will copy all images from the sourcePath to the destinationPath, 
creating a folder structure based on the date the image was taken. The destinationPath 
will need to be created before running the script, inside that folder the script will
