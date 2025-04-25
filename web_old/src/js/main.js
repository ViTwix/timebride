// Main JavaScript entry point
document.addEventListener('DOMContentLoaded', () => {
  // Initialize components
  initializeNavigation();
  initializeFormValidation();
  initializeImagePreviews();
});

function initializeNavigation() {
  const mobileMenuButton = document.querySelector('[data-mobile-menu-button]');
  const mobileMenu = document.querySelector('[data-mobile-menu]');

  if (!mobileMenuButton || !mobileMenu) return;

  mobileMenuButton.addEventListener('click', () => {
    const isExpanded = mobileMenuButton.getAttribute('aria-expanded') === 'true';
    mobileMenuButton.setAttribute('aria-expanded', !isExpanded);
    mobileMenu.classList.toggle('hidden');
  });
}

function initializeFormValidation() {
  const forms = document.querySelectorAll('form[data-validate]');
  
  forms.forEach(form => {
    form.addEventListener('submit', (e) => {
      let isValid = true;
      const requiredFields = form.querySelectorAll('[required]');

      requiredFields.forEach(field => {
        if (!validateField(field)) {
          isValid = false;
        }
      });

      if (!isValid) {
        e.preventDefault();
      }
    });

    // Real-time validation on input
    const fields = form.querySelectorAll('input, textarea, select');
    fields.forEach(field => {
      field.addEventListener('blur', () => {
        validateField(field);
      });

      field.addEventListener('input', () => {
        // Remove error state while typing
        field.classList.remove('error');
        const errorElement = field.parentElement.querySelector('.error-message');
        if (errorElement) {
          errorElement.remove();
        }
      });
    });
  });
}

function validateField(field) {
  const value = field.value.trim();
  let isValid = true;
  let errorMessage = '';

  // Remove existing error message
  const existingError = field.parentElement.querySelector('.error-message');
  if (existingError) {
    existingError.remove();
  }

  // Required field validation
  if (field.hasAttribute('required') && !value) {
    isValid = false;
    errorMessage = 'This field is required';
  }

  // Email validation
  if (field.type === 'email' && value) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(value)) {
      isValid = false;
      errorMessage = 'Please enter a valid email address';
    }
  }

  // File type validation
  if (field.type === 'file' && field.hasAttribute('accept')) {
    const acceptedTypes = field.getAttribute('accept').split(',');
    const file = field.files[0];
    
    if (file && !acceptedTypes.some(type => {
      if (type.startsWith('.')) {
        return file.name.toLowerCase().endsWith(type.toLowerCase());
      }
      return file.type.match(new RegExp(type.replace('*', '.*')));
    })) {
      isValid = false;
      errorMessage = `Please select a valid file type (${acceptedTypes.join(', ')})`;
    }
  }

  // Update UI based on validation
  if (!isValid) {
    field.classList.add('error');
    const errorElement = document.createElement('div');
    errorElement.className = 'error-message text-sm text-red-600 mt-1';
    errorElement.textContent = errorMessage;
    field.parentElement.appendChild(errorElement);
  } else {
    field.classList.remove('error');
  }

  return isValid;
}

function initializeImagePreviews() {
  const imageInputs = document.querySelectorAll('input[type="file"][data-preview]');
  
  imageInputs.forEach(input => {
    const previewId = input.getAttribute('data-preview');
    const previewElement = document.getElementById(previewId);

    if (!previewElement) return;

    input.addEventListener('change', () => {
      const file = input.files[0];
      
      if (file && file.type.startsWith('image/')) {
        const reader = new FileReader();
        
        reader.onload = (e) => {
          previewElement.src = e.target.result;
          previewElement.classList.remove('hidden');
        };
        
        reader.readAsDataURL(file);
      } else {
        previewElement.src = '';
        previewElement.classList.add('hidden');
      }
    });
  });
}

// Utility function to format file size
function formatFileSize(bytes) {
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
} 