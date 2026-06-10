document.addEventListener('DOMContentLoaded', function() {
    const timeline = document.getElementById('adjustList');
    const modal = document.getElementById('modal');
    const modalTitle = document.getElementById('modalTitle');
    const modalBody = document.getElementById('modalBody');
    const closeBtn = document.querySelector('.close');

    fetch('/view/adjust', { method: 'POST' })
    .then(r => r.json())
    .then(data => {
        if (!data || data.length === 0) {
            timeline.innerHTML = '<div class="empty-state"><div class="empty-state-icon">📋</div><p>暂无调整信息</p></div>';
            return;
        }

        data.forEach(adjust => {
            const item = document.createElement('div');
            item.className = 'timeline-item';
            item.innerHTML = `
                <div class="timeline-date">${adjust.adjust_date}</div>
                <div class="timeline-title">${adjust.adjust_title}</div>
                ${adjust.adjust_isbn ? `<span class="timeline-isbn">ISBN: ${adjust.adjust_isbn}</span>` : ''}
                <div class="timeline-preview">${adjust.adjust_content}</div>
            `;
            item.addEventListener('click', () => {
                modalTitle.textContent = adjust.adjust_title;
                modalBody.textContent = adjust.adjust_content;
                modal.classList.add('active');
            });
            timeline.appendChild(item);
        });
    })
    .catch(() => {
        timeline.innerHTML = '<div class="empty-state"><div class="empty-state-icon">⚠️</div><p>加载失败</p></div>';
    });

    closeBtn.onclick = () => modal.classList.remove('active');
    window.onclick = (e) => { if (e.target === modal) modal.classList.remove('active'); };

    // 认证按钮
    const btn = document.getElementById('authButton');
    if (btn) {
        if (localStorage.getItem('isLoggedIn') === 'true') {
            btn.textContent = '个人中心';
            btn.href = '/user/library';
        } else {
            btn.href = '/login?redirect=' + encodeURIComponent(window.location.pathname);
        }
    }
});