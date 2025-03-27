// Admin panel JavaScript
document.addEventListener('DOMContentLoaded', function() {
    
    // Sidebar mobil görünüm kontrolü
    const sidebarToggle = document.getElementById('sidebarToggle');
    const adminSidebar = document.querySelector('.admin-sidebar');
    
    if (sidebarToggle) {
        sidebarToggle.addEventListener('click', function() {
            adminSidebar.classList.toggle('show');
        });
    }
    
    // Form doğrulama
    const forms = document.querySelectorAll('.needs-validation');
    Array.from(forms).forEach(form => {
        form.addEventListener('submit', event => {
            if (!form.checkValidity()) {
                event.preventDefault();
                event.stopPropagation();
            }
            form.classList.add('was-validated');
        }, false);
    });
    
    // Silme işlemi onay modalı
    const confirmDeleteModal = document.getElementById('confirmDeleteModal');
    if (confirmDeleteModal) {
        confirmDeleteModal.addEventListener('show.bs.modal', function(event) {
            const button = event.relatedTarget;
            const id = button.getAttribute('data-id');
            const name = button.getAttribute('data-name');
            const type = button.getAttribute('data-type');
            
            const modalTitle = confirmDeleteModal.querySelector('.modal-title');
            const modalBody = confirmDeleteModal.querySelector('.modal-body p');
            const confirmButton = confirmDeleteModal.querySelector('#confirmDeleteButton');
            
            modalTitle.textContent = `${type} Sil`;
            modalBody.textContent = `"${name}" ${type.toLowerCase()}ını silmek istediğinizden emin misiniz? Bu işlem geri alınamaz.`;
            confirmButton.setAttribute('data-id', id);
        });
        
        const confirmDeleteButton = document.getElementById('confirmDeleteButton');
        if (confirmDeleteButton) {
            confirmDeleteButton.addEventListener('click', function() {
                const id = this.getAttribute('data-id');
                const form = document.getElementById('deleteForm');
                const idInput = document.getElementById('deleteItemId');
                
                idInput.value = id;
                form.submit();
            });
        }
    }
    
    // Medya seçici
    const mediaSelectButtons = document.querySelectorAll('.media-select-button');
    if (mediaSelectButtons.length > 0) {
        mediaSelectButtons.forEach(button => {
            button.addEventListener('click', function() {
                const targetInput = this.getAttribute('data-target');
                const targetPreview = this.getAttribute('data-preview');
                
                // Medya manager modalını aç
                const mediaModal = new bootstrap.Modal(document.getElementById('mediaManagerModal'));
                mediaModal.show();
                
                // Seçilen medya için event listener
                const mediaItems = document.querySelectorAll('.media-item');
                mediaItems.forEach(item => {
                    item.addEventListener('click', function() {
                        const mediaId = this.getAttribute('data-id');
                        const mediaUrl = this.getAttribute('data-url');
                        
                        // Hedef inputa ID'yi ata
                        document.getElementById(targetInput).value = mediaId;
                        
                        // Eğer önizleme varsa URL'yi ata
                        if (targetPreview) {
                            document.getElementById(targetPreview).src = mediaUrl;
                            document.getElementById(targetPreview).classList.remove('d-none');
                        }
                        
                        // Modalı kapat
                        mediaModal.hide();
                    });
                });
            });
        });
    }
    
    // Kategori ağacı düzenleme
    const categoryTree = document.getElementById('categoryTree');
    if (categoryTree) {
        // Dragula veya benzeri bir kütüphane ile sürükle bırak özelliği eklenebilir
    }
    
    // Slug oluşturma (URL dostu metin)
    const titleInput = document.getElementById('title');
    const slugInput = document.getElementById('slug');
    
    if (titleInput && slugInput) {
        titleInput.addEventListener('keyup', function() {
            // Slug alanı boşsa veya başlangıçta oluşturulmuşsa otomatik slug oluştur
            if (slugInput.getAttribute('data-autogenerate') === 'true') {
                slugInput.value = createSlug(this.value);
            }
        });
        
        // Manual slug düzenlendiğinde, otomatik oluşturmayı devre dışı bırak
        slugInput.addEventListener('input', function() {
            this.setAttribute('data-autogenerate', 'false');
        });
        
        // Slug oluşturucu yardımcı fonksiyonu
        function createSlug(text) {
            // Türkçe karakterleri dönüştür
            const turkishChars = {'ç':'c', 'ğ':'g', 'ı':'i', 'ö':'o', 'ş':'s', 'ü':'u', 'Ç':'C', 'Ğ':'G', 'İ':'I', 'Ö':'O', 'Ş':'S', 'Ü':'U'};
            
            return text.toString().toLowerCase()
                .replace(/[çğıöşüÇĞİÖŞÜ]/g, function(letter) { return turkishChars[letter] || letter; })
                .replace(/\s+/g, '-')           // Boşlukları tire ile değiştir
                .replace(/[^\w\-]+/g, '')       // Alfanümerik ve tire dışındaki karakterleri kaldır
                .replace(/\-\-+/g, '-')         // Birden fazla tireyi tek tireye indirge
                .replace(/^-+/, '')             // Baştaki tireleri kaldır
                .replace(/-+$/, '');            // Sondaki tireleri kaldır
        }
    }
    
    // Tarih formatı için flatpickr veya benzeri bir kütüphane kullanılabilir
    
    // Dashboard için grafikler (Chart.js veya benzeri bir kütüphane ile)
    const visitorChart = document.getElementById('visitorChart');
    if (visitorChart) {
        // Chart.js veya benzeri bir kütüphane ile ziyaretçi grafiği oluşturulabilir
    }
}); 