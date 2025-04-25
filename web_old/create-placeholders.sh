#!/bin/bash

# Create placeholders for missing images
mkdir -p public/img/features
mkdir -p public/img/avatars
mkdir -p public/fonts

# Function to create a placeholder SVG
create_svg_placeholder() {
  local output_file="$1"
  local width="$2"
  local height="$3"
  local text="$4"
  
  echo "Creating placeholder: $output_file"
  
  cat > "$output_file" << EOF
<svg width="$width" height="$height" xmlns="http://www.w3.org/2000/svg">
  <rect width="$width" height="$height" fill="#E3D5CA"/>
  <text x="$(($width/2))" y="$(($height/2))" font-family="Arial" font-size="24" fill="#6D4C41" text-anchor="middle">$text</text>
</svg>
EOF
}

# Create feature image placeholders
create_svg_placeholder "public/img/features/calendar-feature.svg" 800 500 "Календар"
create_svg_placeholder "public/img/features/client-feature.svg" 800 500 "Клієнти"
create_svg_placeholder "public/img/features/contract-feature.svg" 800 500 "Документи"
create_svg_placeholder "public/img/photography-bg.svg" 1920 1080 "Фон фотографії"

# Create avatar placeholders
create_svg_placeholder "public/img/avatars/avatar-1.svg" 96 96 "Фото 1"
create_svg_placeholder "public/img/avatars/avatar-2.svg" 96 96 "Фото 2"
create_svg_placeholder "public/img/avatars/avatar-3.svg" 96 96 "Фото 3"

# Create symlinks from SVG to PNG for better browser compatibility
ln -sf calendar-feature.svg public/img/features/calendar-feature.png
ln -sf client-feature.svg public/img/features/client-feature.png
ln -sf contract-feature.svg public/img/features/contract-feature.png
ln -sf ../photography-bg.svg public/img/photography-bg.jpg

# Create symbolic links to images in the main public directory
ln -sf features/calendar-feature.png public/img/calendar-feature.png
ln -sf features/client-feature.png public/img/client-feature.png
ln -sf features/contract-feature.png public/img/contract-feature.png

# Create empty font files
touch public/fonts/Inter-Regular.woff
touch public/fonts/Inter-Regular.woff2
touch public/fonts/Inter-Medium.woff
touch public/fonts/Inter-Medium.woff2
touch public/fonts/Inter-Bold.woff
touch public/fonts/Inter-Bold.woff2
touch public/fonts/Inter-Italic.woff
touch public/fonts/Inter-Italic.woff2

echo "All placeholders created successfully!" 