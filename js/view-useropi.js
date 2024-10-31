document.addEventListener('DOMContentLoaded', function() {
    const searchForm = document.getElementById('searchForm');
    const opinionsList = document.getElementById('opinionsList');
    const replyModal = document.getElementById('replyModal');
    const closeModal = document.querySelector('.close');
    const replyForm = document.getElementById('replyForm');

    // 页面加载时显示所有用户意见
    fetchOpinions('');

    // 搜索表单提交事件
    searchForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const searchId = document.getElementById('searchId').value;
        fetchOpinions(searchId);
    });

    // 获取用户意见
    function fetchOpinions(opinionId) {
        fetch('/view/useropi', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `opinion_id=${opinionId}`
        })
        .then(response => response.json())
        .then(data => {
            displayOpinions(data);
        })
        .catch(error => {
            console.error('Error:', error);
            // alert('获取用户意见失败，请重试。');
        });
    }

    // 显示用户意见
    function displayOpinions(opinions) {
        opinionsList.innerHTML = '';
        opinions.forEach(opinion => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${opinion.opinion_id}</td>
                <td>${opinion.name}</td>
                <td>${opinion.phone}</td>
                <td>${opinion.email}</td>
                <td>${opinion.idea}</td>
                <td><button class="reply-btn" data-id="${opinion.opinion_id}">回复</button></td>
            `;
            opinionsList.appendChild(row);
        });

        // 为所有回复按钮添加事件监听器
        document.querySelectorAll('.reply-btn').forEach(button => {
            button.addEventListener('click', function() {
                openReplyModal(this.getAttribute('data-id'));
            });
        });
    }

    // 打开回复模态框
    function openReplyModal(opinionId) {
        document.getElementById('replay_user').value = opinionId;
        replyModal.style.display = 'block';
    }

    // 关闭模态框
    closeModal.onclick = function() {
        replyModal.style.display = 'none';
    }

    // 点击模态框外部关闭
    window.onclick = function(event) {
        if (event.target == replyModal) {
            replyModal.style.display = 'none';
        }
    }

    // 提交回复
    replyForm.addEventListener('submit', function(e) {
        e.preventDefault();
        if (confirm('确认提交回复？')) {
            const formData = new FormData(replyForm);

            fetch('/replay/useropi', {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (response.ok) {
                    alert('回复提交成功');
                    replyModal.style.display = 'none';
                    replyForm.reset();
                    fetchOpinions(''); // 重新加载所有意见
                } else if (response.status === 400){
                    alert("请当天回复")
            }else {
                    throw new Error('回复提交失败');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('回复提交失败，请重试。');
            });
        }
    });

    // 设置回复日期为今天
    document.getElementById('replay_date').valueAsDate = new Date();
    document.getElementById('replay_name').value = "智慧图书馆"
});