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
    
    // Tüm resim URL'lerini güncelleyip önbelleği kırmak için zaman damgası ekleyen fonksiyon
    const updateImageUrls = () => {
        const allImages = document.querySelectorAll('img[src*="unsplash.com"]');
        
        allImages.forEach(img => {
            // Mevcut URL'i al
            let currentUrl = img.getAttribute('src');
            
            // URL'de zaten ? varsa & ile ekle, yoksa ? ekle
            const separator = currentUrl.includes('?') ? '&' : '?';
            
            // Zaman damgası ekle
            const timestamp = new Date().getTime();
            const newUrl = `${currentUrl}${separator}t=${timestamp}`;
            
            // Yeni URL'i ayarla
            img.setAttribute('src', newUrl);
        });
    };
    
    // Fonksiyonu çağır
    updateImageUrls();
    
    // Varsayılan görsel URL'lerini güncel olanlarla değiştir
    const replaceImageUrls = () => {
        const unsplashReplacements = {
            "istanbulda-su-kesintisi": "https://images.unsplash.com/photo-1563543054715-4f60d55bccb3?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
            "yeni-vergi-duzenlemesi": "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
            "meteorolojiden-kuvvetli-yagis": "https://images.unsplash.com/photo-1514632595-4944383f2737?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
            "dolar-rekor-kirdi": "https://images.unsplash.com/photo-1591696205602-2f950c417cb9?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
            "merkez-bankasi-faiz": "https://images.unsplash.com/photo-1553729459-efe14ef6055d?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
            "asgari-ucret-zammi": "https://images.unsplash.com/photo-1607863680198-23d4b2565df0?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
            "galatasaray-fenerbahce": "https://images.unsplash.com/photo-1508098682722-e99c643e7f3b?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
            "milli-basketbolcu": "https://images.unsplash.com/photo-1518407613690-d9fc990e795f?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
            "turkiye-dunya-kupasi": "https://images.unsplash.com/photo-1522778119026-d647f0596c20?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80"
        };

        const allImages = document.querySelectorAll('img');
        
        allImages.forEach(img => {
            const src = img.getAttribute('src');
            const parentLink = img.closest('a');
            let slug = "";
            
            // Eğer img bir a elementi içindeyse, a elementinin href'indeki slug'ı al
            if (parentLink) {
                const href = parentLink.getAttribute('href');
                if (href && href.includes('/haber/')) {
                    slug = href.split('/haber/')[1];
                }
            }
            
            // Resim URL'si yoksa veya boşsa ya da placeholder içeriyorsa
            if (!src || src === '' || src.includes('placeholder')) {
                // Slug'dan anahtar kelime kontrolü yap
                for (const [key, url] of Object.entries(unsplashReplacements)) {
                    if (slug.includes(key)) {
                        img.setAttribute('src', url);
                        break;
                    }
                }
            }
            
            // Resim Unsplash URL'si değilse, doğrudan slug'a göre eşleşen bir resim ile değiştir
            if (src && !src.includes('unsplash.com')) {
                for (const [key, url] of Object.entries(unsplashReplacements)) {
                    if (slug.includes(key)) {
                        img.setAttribute('src', url);
                        break;
                    }
                }
                
                // Eğer hiçbir slug eşleşmediyse ve resim görüntülenemiyorsa veya boyutu çok küçükse varsayılan bir resim ata
                img.onerror = function() {
                    this.src = "https://images.unsplash.com/photo-1586339949916-3e9457bef6d3?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80";
                    this.onerror = null; // Sonsuz döngüyü önle
                };
            }
        });
    };

    // Sayfa yüklendiğinde replaceImageUrls fonksiyonunu çağır
    replaceImageUrls();
}); 