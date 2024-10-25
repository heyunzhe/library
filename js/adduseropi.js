document.addEventListener('DOMContentLoaded', function() {
    const modal = document.getElementById('userOpiModal');
    const form = document.getElementById('userOpiForm');
    const closeBtn = modal.querySelector('.close2');
    
    document.getElementById('openModalBtn').onclick = openModal;
    // Function to open the modal
    function openModal() {
        modal.style.display = 'block';
    }

    // Function to close the modal
    function closeModal() {
        modal.style.display = 'none';
    }

    // Close the modal when clicking on <span> (x)
    closeBtn.onclick = closeModal;

    // Close the modal when clicking outside of it
    window.onclick = function(event) {
        if (event.target == modal) {
            closeModal();
        }
    }

    form.addEventListener('submit', function(e) {
        e.preventDefault();

        const formData = new FormData(form);

        fetch('/add/useropi', {
            method: 'POST',
            body: formData
        })
        .then(response => {
            if (response.ok) {
                alert('意见提交成功，感谢您的反馈！');
                form.reset();
                closeModal();
            } else if (response.status === 500) {
                alert('服务器错误，请稍后重试。');
            } else {
                alert('提交失败，请重试。');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('提交失败，请检查您的网络连接并重试。');
        });
    });


});