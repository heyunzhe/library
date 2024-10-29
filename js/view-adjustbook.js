document.addEventListener('DOMContentLoaded', function() {
    const adjustList = document.getElementById('adjust-list');
    const modal = document.getElementById('modal');
    const modalTitle = document.getElementById('modal-title');
    const modalBody = document.getElementById('modal-body');
    const closeBtn = document.getElementsByClassName('close')[0];

    // 获取调整信息
    fetch('/view/adjust', {
        method: 'POST'
    })
    .then(response => response.json())
    .then(data => {
        data.forEach(adjust => {
            const adjustItem = document.createElement('div');
            adjustItem.className = 'adjust-item';
            adjustItem.innerHTML = `
                <span class="adjust-date">${adjust.adjust_date}</span>
                <a href="#" class="adjust-title">${adjust.adjust_title}</a>
            `;
            adjustList.appendChild(adjustItem);

            const titleLink = adjustItem.querySelector('.adjust-title');
            titleLink.addEventListener('click', function(e) {
                e.preventDefault();
                showModal(adjust);
            });
        });
    })
    .catch(error => console.error('Error:', error));

    function showModal(adjust) {
        modalTitle.textContent = adjust.adjust_title;
        modalBody.textContent = adjust.adjust_content;
        modal.style.display = 'block';
    }

    closeBtn.onclick = function() {
        modal.style.display = 'none';
    }

    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = 'none';
        }
    }
});