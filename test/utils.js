// 通用工具函数
class Utils {
    // 显示提示信息
    static showAlert(message, type = 'info', container = 'alertsContainer') {   
        const alertDiv = document.createElement('div');
        alertDiv.className = `alert alert-${type} alert-dismissible fade show`;
        alertDiv.innerHTML = `
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        `;
        
        const containerEl = document.getElementById(container);
        if (containerEl) {
            containerEl.appendChild(alertDiv);
            
            // 5秒后自动移除
            setTimeout(() => {
                if (alertDiv.parentNode) {
                    alertDiv.remove();
                }
            }, 5000);
        }
    }

    // 转义HTML
    static escapeHtml(unsafe) {
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }

    // 截断文本
    static truncateText(text, maxLength) {
        if (!text) return '';
        if (text.length <= maxLength) return text;
        return text.substring(0, maxLength) + '...';
    }

    // 格式化时间
    static formatDateTime(dateString) {
        if (!dateString) return 'N/A';
        const date = new Date(dateString);
        return date.toLocaleString('zh-CN');
    }

    // 格式化时长
    static formatDuration(seconds) {
        if (!seconds) return '0s';
        if (seconds < 60) return `${seconds.toFixed(2)}s`;
        const minutes = Math.floor(seconds / 60);
        const remainingSeconds = seconds % 60;
        return `${minutes}m ${remainingSeconds.toFixed(2)}s`;
    }

    // 获取状态徽章类名
    static getStatusBadgeClass(status) {
        const statusClasses = {
            'passed': 'bg-success',
            'failed': 'bg-danger',
            'error': 'bg-warning',
            'running': 'bg-primary',
            'completed': 'bg-success',
            'cancelled': 'bg-secondary',
            'active': 'bg-success',
            'inactive': 'bg-secondary',
            'enabled': 'bg-success',
            'disabled': 'bg-secondary'
        };
        return statusClasses[status] || 'bg-secondary';
    }

    // 获取级别徽章类名
    static getLevelBadgeClass(level) {
        const levelClasses = {
            'DEBUG': 'bg-secondary',
            'INFO': 'bg-info',
            'WARNING': 'bg-warning',
            'ERROR': 'bg-danger',
            'CRITICAL': 'bg-dark'
        };
        return levelClasses[level] || 'bg-secondary';
    }

    // 显示加载状态
    static showLoading(element) {
        if (element) {
            element.classList.add('loading');
        }
    }

    // 隐藏加载状态
    static hideLoading(element) {
        if (element) {
            element.classList.remove('loading');
        }
    }

    // 验证表单
    static validateForm(formElement) {
        const inputs = formElement.querySelectorAll('input[required], select[required], textarea[required]');
        let isValid = true;

        inputs.forEach(input => {
            if (!input.value.trim()) {
                input.classList.add('is-invalid');
                isValid = false;
            } else {
                input.classList.remove('is-invalid');
            }
        });

        return isValid;
    }

    // 重置表单
    static resetForm(formElement) {
        formElement.reset();
        // 清除验证状态
        const inputs = formElement.querySelectorAll('.is-invalid');
        inputs.forEach(input => input.classList.remove('is-invalid'));
    }

    // 复制到剪贴板
    static copyToClipboard(text) {
        navigator.clipboard.writeText(text).then(() => {
            this.showAlert('已复制到剪贴板', 'success');
        }).catch(() => {
            this.showAlert('复制失败', 'error');
        });
    }

    // 防抖函数
    static debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }

    // 生成分页HTML
    static generatePagination(currentPage, totalPages, containerId) {
        const container = document.getElementById(containerId);
        if (!container || totalPages <= 1) {
            container.innerHTML = '';
            return;
        }

        let html = '';
        const maxVisiblePages = 5;
        let startPage = Math.max(1, currentPage - Math.floor(maxVisiblePages / 2));
        let endPage = Math.min(totalPages, startPage + maxVisiblePages - 1);

        if (endPage - startPage + 1 < maxVisiblePages) {
            startPage = Math.max(1, endPage - maxVisiblePages + 1);
        }

        // 上一页
        if (currentPage > 1) {
            html += `<li class="page-item"><a class="page-link" href="#" data-page="${currentPage - 1}">上一页</a></li>`;
        }

        // 第一页
        if (startPage > 1) {
            html += `<li class="page-item"><a class="page-link" href="#" data-page="1">1</a></li>`;
            if (startPage > 2) {
                html += `<li class="page-item disabled"><span class="page-link">...</span></li>`;
            }
        }

        // 页码
        for (let i = startPage; i <= endPage; i++) {
            const active = i === currentPage ? 'active' : '';
            html += `<li class="page-item ${active}"><a class="page-link" href="#" data-page="${i}">${i}</a></li>`;
        }

        // 最后一页
        if (endPage < totalPages) {
            if (endPage < totalPages - 1) {
                html += `<li class="page-item disabled"><span class="page-link">...</span></li>`;
            }
            html += `<li class="page-item"><a class="page-link" href="#" data-page="${totalPages}">${totalPages}</a></li>`;
        }

        // 下一页
        if (currentPage < totalPages) {
            html += `<li class="page-item"><a class="page-link" href="#" data-page="${currentPage + 1}">下一页</a></li>`;
        }

        container.innerHTML = html;
    }
}

// 导出到全局作用域
window.Utils = Utils;