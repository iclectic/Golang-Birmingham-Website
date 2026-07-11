// Golang Birmingham - Progressive Enhancement JavaScript
// This script provides minimal enhancements that work without JavaScript

(function() {
    'use strict';

    // Mobile Navigation Toggle
    function initMobileNavigation() {
        const navToggle = document.querySelector('.nav-toggle');
        const mobileMenu = document.querySelector('.nav-mobile');
        
        if (!navToggle || !mobileMenu) return;

        navToggle.addEventListener('click', function() {
            const isExpanded = navToggle.getAttribute('aria-expanded') === 'true';
            
            navToggle.setAttribute('aria-expanded', !isExpanded);
            mobileMenu.classList.toggle('nav-mobile--open');
            
            // Prevent body scroll when menu is open
            document.body.style.overflow = isExpanded ? '' : 'hidden';
        });

        // Close menu when clicking outside
        document.addEventListener('click', function(event) {
            if (!navToggle.contains(event.target) && !mobileMenu.contains(event.target)) {
                navToggle.setAttribute('aria-expanded', 'false');
                mobileMenu.classList.remove('nav-mobile--open');
                document.body.style.overflow = '';
            }
        });

        // Close menu when pressing Escape
        document.addEventListener('keydown', function(event) {
            if (event.key === 'Escape' && mobileMenu.classList.contains('nav-mobile--open')) {
                navToggle.setAttribute('aria-expanded', 'false');
                mobileMenu.classList.remove('nav-mobile--open');
                document.body.style.overflow = '';
            }
        });

        // Close menu when navigation item is clicked
        const navLinks = mobileMenu.querySelectorAll('.nav-link');
        navLinks.forEach(function(link) {
            link.addEventListener('click', function() {
                navToggle.setAttribute('aria-expanded', 'false');
                mobileMenu.classList.remove('nav-mobile--open');
                document.body.style.overflow = '';
            });
        });
    }

    // Gallery Lightbox
    function initGalleryLightbox() {
        const galleryLinks = document.querySelectorAll('.gallery-link');
        const lightbox = document.getElementById('lightbox');
        
        if (!galleryLinks.length || !lightbox) return;

        const lightboxImage = document.getElementById('lightbox-image');
        const lightboxTitle = document.getElementById('lightbox-title');
        const lightboxCaption = document.getElementById('lightbox-caption');
        const lightboxClose = document.getElementById('lightbox-close');
        const lightboxPrev = document.getElementById('lightbox-prev');
        const lightboxNext = document.getElementById('lightbox-next');

        let currentIndex = 0;
        let galleryItems = Array.from(galleryLinks);

        function openLightbox(index) {
            currentIndex = index;
            const item = galleryItems[index];
            const img = item.querySelector('img');
            const caption = item.querySelector('.gallery-caption');
            const event = item.querySelector('.gallery-event');
            const category = item.querySelector('.gallery-category');

            if (img) {
                // Use full-size image (remove thumbnail suffix)
                const fullImageSrc = img.src.replace('/thumbnails/', '/');
                lightboxImage.src = fullImageSrc;
                lightboxImage.alt = img.alt;
            }

            if (lightboxTitle && event) {
                lightboxTitle.textContent = event.textContent;
            }

            if (lightboxCaption && caption) {
                lightboxCaption.textContent = caption.textContent;
            }

            lightbox.style.display = 'flex';
            document.body.style.overflow = 'hidden';
            
            // Focus management
            lightboxClose.focus();
        }

        function closeLightbox() {
            lightbox.style.display = 'none';
            document.body.style.overflow = '';
        }

        function navigateLightbox(direction) {
            if (direction === 'prev') {
                currentIndex = (currentIndex - 1 + galleryItems.length) % galleryItems.length;
            } else {
                currentIndex = (currentIndex + 1) % galleryItems.length;
            }
            openLightbox(currentIndex);
        }

        // Event listeners
        galleryLinks.forEach(function(link, index) {
            link.addEventListener('click', function(e) {
                e.preventDefault();
                openLightbox(index);
            });
        });

        if (lightboxClose) {
            lightboxClose.addEventListener('click', closeLightbox);
        }

        if (lightboxPrev) {
            lightboxPrev.addEventListener('click', function() {
                navigateLightbox('prev');
            });
        }

        if (lightboxNext) {
            lightboxNext.addEventListener('click', function() {
                navigateLightbox('next');
            });
        }

        // Keyboard navigation
        document.addEventListener('keydown', function(e) {
            if (lightbox.style.display === 'flex') {
                switch (e.key) {
                    case 'Escape':
                        closeLightbox();
                        break;
                    case 'ArrowLeft':
                        navigateLightbox('prev');
                        break;
                    case 'ArrowRight':
                        navigateLightbox('next');
                        break;
                }
            }
        });

        // Close on background click
        lightbox.addEventListener('click', function(e) {
            if (e.target === lightbox) {
                closeLightbox();
            }
        });
    }

    // Gallery Filtering
    function initGalleryFilters() {
        const eventFilter = document.getElementById('event-filter');
        const categoryFilter = document.getElementById('category-filter');
        const yearFilter = document.getElementById('year-filter');
        const clearFilters = document.getElementById('clear-filters');
        const galleryItems = document.querySelectorAll('.gallery-item');

        if (!galleryItems.length) return;

        function filterGallery() {
            const eventValue = eventFilter ? eventFilter.value : '';
            const categoryValue = categoryFilter ? categoryFilter.value : '';
            const yearValue = yearFilter ? yearFilter.value : '';

            galleryItems.forEach(function(item) {
                const itemEvent = item.dataset.event || '';
                const itemCategory = item.dataset.category || '';
                const itemYear = item.dataset.year || '';

                const matchesEvent = !eventValue || itemEvent === eventValue;
                const matchesCategory = !categoryValue || itemCategory === categoryValue;
                const matchesYear = !yearValue || itemYear === yearValue;

                if (matchesEvent && matchesCategory && matchesYear) {
                    item.style.display = '';
                } else {
                    item.style.display = 'none';
                }
            });
        }

        if (eventFilter) {
            eventFilter.addEventListener('change', filterGallery);
        }

        if (categoryFilter) {
            categoryFilter.addEventListener('change', filterGallery);
        }

        if (yearFilter) {
            yearFilter.addEventListener('change', filterGallery);
        }

        if (clearFilters) {
            clearFilters.addEventListener('click', function() {
                if (eventFilter) eventFilter.value = '';
                if (categoryFilter) categoryFilter.value = '';
                if (yearFilter) yearFilter.value = '';
                filterGallery();
            });
        }
    }

    // Form Enhancements
    function initFormEnhancements() {
        const forms = document.querySelectorAll('form');
        
        forms.forEach(function(form) {
            // Add client-side validation feedback
            const requiredFields = form.querySelectorAll('[required]');
            
            requiredFields.forEach(function(field) {
                field.addEventListener('blur', function() {
                    if (!field.value.trim()) {
                        field.classList.add('error');
                    } else {
                        field.classList.remove('error');
                    }
                });

                field.addEventListener('input', function() {
                    if (field.value.trim()) {
                        field.classList.remove('error');
                    }
                });
            });

            // Form submission enhancement
            form.addEventListener('submit', function(e) {
                let isValid = true;
                
                requiredFields.forEach(function(field) {
                    if (!field.value.trim()) {
                        field.classList.add('error');
                        isValid = false;
                    }
                });

                if (!isValid) {
                    e.preventDefault();
                    // Focus first error field
                    const firstError = form.querySelector('.error');
                    if (firstError) {
                        firstError.focus();
                    }
                }
            });
        });
    }

    // Smooth Scrolling for anchor links
    function initSmoothScrolling() {
        const anchorLinks = document.querySelectorAll('a[href^="#"]');
        
        anchorLinks.forEach(function(link) {
            link.addEventListener('click', function(e) {
                const targetId = link.getAttribute('href').substring(1);
                const targetElement = document.getElementById(targetId);
                
                if (targetElement) {
                    e.preventDefault();
                    targetElement.scrollIntoView({
                        behavior: 'smooth',
                        block: 'start'
                    });
                }
            });
        });
    }

    // Copy link functionality
    function initCopyLink() {
        const copyButtons = document.querySelectorAll('[data-copy-url]');
        
        copyButtons.forEach(function(button) {
            button.addEventListener('click', function() {
                const url = button.dataset.copyUrl || window.location.href;
                
                if (navigator.clipboard) {
                    navigator.clipboard.writeText(url).then(function() {
                        const originalText = button.textContent;
                        button.textContent = 'Copied!';
                        button.classList.add('copied');
                        
                        setTimeout(function() {
                            button.textContent = originalText;
                            button.classList.remove('copied');
                        }, 2000);
                    });
                }
            });
        });
    }

    // Scroll-based animations (subtle)
    function initScrollAnimations() {
        const animatedElements = document.querySelectorAll('[data-animate]');
        
        if (!animatedElements.length) return;

        function checkScroll() {
            animatedElements.forEach(function(element) {
                const rect = element.getBoundingClientRect();
                const isVisible = rect.top < window.innerHeight && rect.bottom > 0;
                
                if (isVisible) {
                    element.classList.add('animate-in');
                }
            });
        }

        // Initial check
        checkScroll();
        
        // Check on scroll
        window.addEventListener('scroll', checkScroll);
    }

    // Initialize all features when DOM is ready
    function init() {
        initMobileNavigation();
        initGalleryLightbox();
        initGalleryFilters();
        initFormEnhancements();
        initSmoothScrolling();
        initCopyLink();
        initScrollAnimations();
    }

    // Wait for DOM to be ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }

})();
