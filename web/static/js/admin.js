/**
 * Admin Panel JavaScript
 */

// DOM yüklendikten sonra çalışacak kodlar
document.addEventListener("DOMContentLoaded", function() {
    // Sidebar toggle butonu
    const sidebarToggleBtn = document.querySelector('.sidebar-toggle-btn');
    if (sidebarToggleBtn) {
        sidebarToggleBtn.addEventListener('click', function() {
            document.querySelector('.admin-layout').classList.toggle('sidebar-collapsed');
        });
    }

    // Mobilde sidebar kapatma
    const mediaQuery = window.matchMedia('(max-width: 992px)');
    if (mediaQuery.matches) {
        document.querySelectorAll('.admin-sidebar .nav-link').forEach(function(link) {
            link.addEventListener('click', function() {
                if (window.innerWidth <= 576) {
                    document.querySelector('.admin-layout').classList.add('sidebar-collapsed');
                }
            });
        });
    }

    // Tablo satır seçimi
    const selectableRows = document.querySelectorAll('.selectable-row');
    if (selectableRows.length > 0) {
        selectableRows.forEach(function(row) {
            row.addEventListener('click', function(e) {
                // Butonlara tıklandığında propagation engellenmiş olmalı
                if (!e.target.closest('button') && !e.target.closest('a') && !e.target.closest('input')) {
                    this.classList.toggle('selected');
                    const checkbox = this.querySelector('input[type="checkbox"]');
                    if (checkbox) {
                        checkbox.checked = !checkbox.checked;
                        updateBulkActions();
                    }
                }
            });
        });
    }

    // Tüm checkbox'ları seçme
    const selectAllCheckbox = document.querySelector('#selectAll');
    if (selectAllCheckbox) {
        selectAllCheckbox.addEventListener('change', function() {
            const checkboxes = document.querySelectorAll('.item-checkbox');
            checkboxes.forEach(function(checkbox) {
                checkbox.checked = selectAllCheckbox.checked;
                const row = checkbox.closest('tr');
                if (row) {
                    if (checkbox.checked) {
                        row.classList.add('selected');
                    } else {
                        row.classList.remove('selected');
                    }
                }
            });
            updateBulkActions();
        });
    }

    // Checkbox'lar değiştiğinde satır stilini güncelleme
    const itemCheckboxes = document.querySelectorAll('.item-checkbox');
    if (itemCheckboxes.length > 0) {
        itemCheckboxes.forEach(function(checkbox) {
            checkbox.addEventListener('change', function(e) {
                e.stopPropagation();
                const row = this.closest('tr');
                if (row) {
                    if (this.checked) {
                        row.classList.add('selected');
                    } else {
                        row.classList.remove('selected');
                    }
                }
                updateBulkActions();
            });
        });
    }

    // Toplu işlemler butonunu güncelleme
    function updateBulkActions() {
        const bulkActionBtn = document.querySelector('.bulk-actions-btn');
        const selectedCount = document.querySelectorAll('.item-checkbox:checked').length;
        
        if (bulkActionBtn) {
            if (selectedCount > 0) {
                bulkActionBtn.classList.remove('disabled');
                bulkActionBtn.querySelector('.count').textContent = selectedCount;
            } else {
                bulkActionBtn.classList.add('disabled');
                bulkActionBtn.querySelector('.count').textContent = '0';
            }
        }
    }

    // Form doğrulama
    const forms = document.querySelectorAll('.needs-validation');
    if (forms.length > 0) {
        forms.forEach(function(form) {
            form.addEventListener('submit', function(event) {
                if (!form.checkValidity()) {
                    event.preventDefault();
                    event.stopPropagation();
                }
                form.classList.add('was-validated');
            }, false);
        });
    }

    // Bildirim kapatma
    const notifications = document.querySelectorAll('.admin-notification');
    if (notifications.length > 0) {
        notifications.forEach(function(notification) {
            const closeBtn = notification.querySelector('.close-btn');
            if (closeBtn) {
                closeBtn.addEventListener('click', function() {
                    notification.classList.add('closing');
                    setTimeout(function() {
                        notification.remove();
                    }, 300);
                });
            }

            // Otomatik kapanma
            if (notification.dataset.autoDismiss) {
                const delay = parseInt(notification.dataset.autoDismiss) || 5000;
                setTimeout(function() {
                    notification.classList.add('closing');
                    setTimeout(function() {
                        notification.remove();
                    }, 300);
                }, delay);
            }
        });
    }

    // Dosya yükleme önizleme
    const fileInputs = document.querySelectorAll('.custom-file-input');
    if (fileInputs.length > 0) {
        fileInputs.forEach(function(input) {
            input.addEventListener('change', function(e) {
                const fileName = this.files[0].name;
                const previewContainer = document.querySelector(this.dataset.preview);
                const fileLabel = this.nextElementSibling;
                
                if (fileLabel) {
                    fileLabel.textContent = fileName;
                }
                
                if (previewContainer && this.files && this.files[0]) {
                    const reader = new FileReader();
                    
                    reader.onload = function(e) {
                        if (previewContainer.tagName === 'IMG') {
                            previewContainer.src = e.target.result;
                        } else {
                            previewContainer.style.backgroundImage = `url('${e.target.result}')`;
                        }
                    }
                    
                    reader.readAsDataURL(this.files[0]);
                }
            });
        });
    }

    // Tarih formatı
    const formatDate = function(date) {
        const d = new Date(date);
        const day = String(d.getDate()).padStart(2, '0');
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const year = d.getFullYear();
        return `${day}.${month}.${year}`;
    };

    // Responsive tablo için data-label ekleme
    const tables = document.querySelectorAll('.table-responsive-stack');
    if (tables.length > 0) {
        tables.forEach(function(table) {
            const thArray = [];
            table.querySelectorAll('thead th').forEach(function(th) {
                thArray.push(th.textContent);
            });
            
            table.querySelectorAll('tbody tr').forEach(function(tr) {
                tr.querySelectorAll('td').forEach(function(td, index) {
                    td.setAttribute('data-label', thArray[index]);
                });
            });
        });
    }

    // Modal içinde form submit
    const modalForms = document.querySelectorAll('.modal form');
    if (modalForms.length > 0) {
        modalForms.forEach(function(form) {
            form.addEventListener('submit', function(e) {
                e.preventDefault();
                
                const formData = new FormData(form);
                const url = form.getAttribute('action');
                const method = form.getAttribute('method') || 'POST';
                
                fetch(url, {
                    method: method,
                    body: formData
                })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        // Başarılı olduğunda
                        const modal = bootstrap.Modal.getInstance(form.closest('.modal'));
                        modal.hide();
                        
                        // Sayfayı yenile veya bildirim göster
                        if (data.redirect) {
                            window.location.href = data.redirect;
                        } else {
                            showNotification('Başarılı', data.message || 'İşlem başarıyla tamamlandı', 'success');
                            if (data.reload) {
                                setTimeout(() => {
                                    window.location.reload();
                                }, 1000);
                            }
                        }
                    } else {
                        // Hata durumunda
                        showNotification('Hata', data.message || 'Bir hata oluştu', 'error');
                    }
                })
                .catch(error => {
                    showNotification('Hata', 'Bir hata oluştu: ' + error.message, 'error');
                });
            });
        });
    }

    // Bildirim gösterme
    window.showNotification = function(title, message, type = 'info') {
        const container = document.querySelector('.notification-container');
        if (!container) return;
        
        const notification = document.createElement('div');
        notification.className = `admin-notification ${type}`;
        notification.dataset.autoDismiss = '5000';
        
        notification.innerHTML = `
            <div class="notification-icon">
                <i class="bi bi-${type === 'success' ? 'check-circle' : type === 'error' ? 'exclamation-circle' : 'info-circle'}-fill"></i>
            </div>
            <div class="notification-content">
                <h6>${title}</h6>
                <p>${message}</p>
            </div>
            <button class="close-btn">
                <i class="bi bi-x"></i>
            </button>
        `;
        
        container.appendChild(notification);
        
        // Animasyon için zaman vermek
        setTimeout(() => {
            notification.classList.add('show');
        }, 10);
        
        // Kapatma butonu
        notification.querySelector('.close-btn').addEventListener('click', function() {
            notification.classList.remove('show');
            setTimeout(() => {
                notification.remove();
            }, 300);
        });
        
        // Otomatik kapanma
        setTimeout(() => {
            notification.classList.remove('show');
            setTimeout(() => {
                notification.remove();
            }, 300);
        }, 5000);
    };

    // Sayfa işlemleri için AJAX fonksiyonları
    window.adminAPI = {
        // İçerik silme
        deleteItem: function(url, itemId, confirmMessage = 'Bu öğeyi silmek istediğinize emin misiniz?') {
            if (confirm(confirmMessage)) {
                fetch(`${url}/${itemId}`, {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        showNotification('Başarılı', data.message || 'Öğe başarıyla silindi', 'success');
                        const elementToRemove = document.querySelector(`[data-id="${itemId}"]`);
                        if (elementToRemove) {
                            elementToRemove.remove();
                        } else {
                            setTimeout(() => {
                                window.location.reload();
                            }, 1000);
                        }
                    } else {
                        showNotification('Hata', data.message || 'Silme işlemi başarısız oldu', 'error');
                    }
                })
                .catch(error => {
                    showNotification('Hata', 'Bir hata oluştu: ' + error.message, 'error');
                });
            }
        },
        
        // Durumu güncelleme
        updateStatus: function(url, itemId, status) {
            fetch(`${url}/${itemId}/status`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ status: status })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showNotification('Başarılı', data.message || 'Durum başarıyla güncellendi', 'success');
                    // Status badge güncelleme
                    const statusBadge = document.querySelector(`[data-id="${itemId}"] .status-badge`);
                    if (statusBadge) {
                        statusBadge.className = `status-badge status-badge-${status.toLowerCase()}`;
                        statusBadge.textContent = status;
                    } else {
                        setTimeout(() => {
                            window.location.reload();
                        }, 1000);
                    }
                } else {
                    showNotification('Hata', data.message || 'Durum güncelleme başarısız oldu', 'error');
                }
            })
            .catch(error => {
                showNotification('Hata', 'Bir hata oluştu: ' + error.message, 'error');
            });
        },
        
        // Toplu işlemler
        bulkAction: function(url, action, ids) {
            if (!ids || ids.length === 0) {
                showNotification('Uyarı', 'Lütfen en az bir öğe seçin', 'warning');
                return;
            }
            
            fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    action: action,
                    ids: ids
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showNotification('Başarılı', data.message || 'İşlem başarıyla tamamlandı', 'success');
                    setTimeout(() => {
                        window.location.reload();
                    }, 1000);
                } else {
                    showNotification('Hata', data.message || 'İşlem başarısız oldu', 'error');
                }
            })
            .catch(error => {
                showNotification('Hata', 'Bir hata oluştu: ' + error.message, 'error');
            });
        },
        
        // Sıralama güncelleme
        updateOrder: function(url, items) {
            fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ items: items })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showNotification('Başarılı', data.message || 'Sıralama başarıyla güncellendi', 'success');
                } else {
                    showNotification('Hata', data.message || 'Sıralama güncelleme başarısız oldu', 'error');
                }
            })
            .catch(error => {
                showNotification('Hata', 'Bir hata oluştu: ' + error.message, 'error');
            });
        }
    };
}); 