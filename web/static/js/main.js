// Ana JavaScript dosyası
document.addEventListener('DOMContentLoaded', function() {
    
    // Bootstrap Tooltip'leri etkinleştir
    var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
    tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl);
    });
    
    // Arama formu doğrulama
    const searchForm = document.querySelector('form[action="/arama"]');
    if (searchForm) {
        searchForm.addEventListener('submit', function(e) {
            const searchInput = this.querySelector('input[name="q"]');
            if (!searchInput.value.trim()) {
                e.preventDefault();
                searchInput.classList.add('is-invalid');
            } else {
                searchInput.classList.remove('is-invalid');
            }
        });
    }
    
    // Yorum formu doğrulama
    const commentForm = document.getElementById('commentForm');
    if (commentForm) {
        commentForm.addEventListener('submit', function(e) {
            const commentContent = document.getElementById('commentContent');
            if (!commentContent.value.trim()) {
                e.preventDefault();
                commentContent.classList.add('is-invalid');
            } else {
                commentContent.classList.remove('is-invalid');
            }
        });
    }
    
    // Lazy load için görünürlük kontrolü
    const lazyImages = document.querySelectorAll('img[data-src]');
    if (lazyImages.length > 0) {
        const imageObserver = new IntersectionObserver(function(entries, observer) {
            entries.forEach(function(entry) {
                if (entry.isIntersecting) {
                    const img = entry.target;
                    img.src = img.dataset.src;
                    img.removeAttribute('data-src');
                    imageObserver.unobserve(img);
                }
            });
        });
        
        lazyImages.forEach(function(img) {
            imageObserver.observe(img);
        });
    }
    
    // Gece modu geçişi
    const darkModeToggle = document.getElementById('darkModeToggle');
    if (darkModeToggle) {
        const prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)');
        
        // Local storage'dan kullanıcı tercihi kontrolü
        const currentTheme = localStorage.getItem('theme');
        if (currentTheme === 'dark') {
            document.body.classList.add('dark-mode');
            darkModeToggle.checked = true;
        } else if (currentTheme === 'light') {
            document.body.classList.remove('dark-mode');
            darkModeToggle.checked = false;
        }
        
        darkModeToggle.addEventListener('change', function() {
            if (this.checked) {
                document.body.classList.add('dark-mode');
                localStorage.setItem('theme', 'dark');
            } else {
                document.body.classList.remove('dark-mode');
                localStorage.setItem('theme', 'light');
            }
        });
    }
    
    // Video iframe'lerini responsive yap
    const articleContent = document.querySelector('.article-content');
    if (articleContent) {
        const iframes = articleContent.querySelectorAll('iframe');
        iframes.forEach(function(iframe) {
            const wrapper = document.createElement('div');
            wrapper.classList.add('ratio', 'ratio-16x9', 'my-4');
            iframe.parentNode.insertBefore(wrapper, iframe);
            wrapper.appendChild(iframe);
        });
    }
    
    // Haberleri kategoriye göre filtreleme
    const categoryFilters = document.querySelectorAll('.category-filter');
    if (categoryFilters.length > 0) {
        categoryFilters.forEach(function(filter) {
            filter.addEventListener('click', function(e) {
                e.preventDefault();
                
                // Aktif filtreyi güncelle
                categoryFilters.forEach(btn => btn.classList.remove('active'));
                this.classList.add('active');
                
                const category = this.dataset.category;
                const articles = document.querySelectorAll('.article-item');
                
                articles.forEach(function(article) {
                    if (category === 'all' || article.dataset.category === category) {
                        article.style.display = 'block';
                    } else {
                        article.style.display = 'none';
                    }
                });
            });
        });
    }
}); 