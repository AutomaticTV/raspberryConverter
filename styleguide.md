# Style guide
This document will reference how to make changes on the aesthetic parts of the project

## Images
The project has basically three different images, in order to change any of the image just replace the existing ones, make sure they have same name and resolution:

### Player background
This image will be displayed through the HDMI port of the Raspberry while there is no content being played.

It's important to note that some text (IP address / LOADING / NO INTERNET) will be added on top of the image. The position, size and color of this text can be adjusted at player/imageMaker.go

* Filepath:  services/player/assets/bg.png
* Resolution: 1920 x 1080

### WEB UI background
This image is used in all the different pages of the web as a background.

* Filepath:  frontend/static/img/bg.jpg
* Resolution: any (minimum recommended 960 x 365)

### WEB UI icon
This image is used as favicon for the web (icon displayed on the tab of the browser)

* Filepath:  frontend/static/img/favicon*.png
* Resolution: 279 x 318

## Others

### Player font
Font used to put text on top of the player background. Replacce for a .tff font file

* Filepath:  services/player/assets/font.tff

### Web style
The web is styled using css. The file containing the css is: frontend/static/main.css
