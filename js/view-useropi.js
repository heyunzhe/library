document.addEventListener('DOMContentLoaded', function() {

    const searchForm = document.getElementById('searchForm');
    const opinionsList = document.getElementById('opinionsList');

    const replyModal = document.getElementById('replyModal');
    const closeModal = document.querySelector('.close');
    const replyForm = document.getElementById('replyForm');

    const viewModal = document.getElementById('viewModal');
    const closeView = document.querySelector('.close-view');

    // 页面加载
    fetchOpinions('');

    // 搜索
    searchForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const searchId = document.getElementById('searchId').value;
        fetchOpinions(searchId);
    });

    // 获取数据
    function fetchOpinions(opinionId) {
        fetch('/view/useropi', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `opinion_id=${opinionId}`
        })
        .then(res => res.json())
        .then(data => displayOpinions(data))
        .catch(err => console.error(err));
    }

    // 渲染表格（关键优化）
    function displayOpinions(opinions) {
        opinionsList.innerHTML = '';

        opinions.forEach(opinion => {

            const shortIdea = opinion.idea.length > 20 
                ? opinion.idea.substring(0, 20) + '...' 
                : opinion.idea;

            const row = document.createElement('tr');

            row.innerHTML = `
                <td>${opinion.opinion_id}</td>
                <td>${opinion.name}</td>
                <td>${opinion.phone}</td>
                <td>${opinion.email}</td>
                <td>${shortIdea}</td>
                <td>
                    <button class="view-btn"
                        data-idea="${encodeURIComponent(opinion.idea)}"
                        data-name="${opinion.name}"
                        data-phone="${opinion.phone}"
                        data-email="${opinion.email}">
                        查看
                    </button>

                    <button class="reply-btn" data-id="${opinion.opinion_id}">
                        回复
                    </button>
                </td>
            `;

            opinionsList.appendChild(row);
        });

        // 查看按钮
        document.querySelectorAll('.view-btn').forEach(btn => {
            btn.addEventListener('click', function() {
                openViewModal(this);
            });
        });

        // 回复按钮
        document.querySelectorAll('.reply-btn').forEach(btn => {
            btn.addEventListener('click', function() {
                openReplyModal(this.getAttribute('data-id'));
            });
        });
    }

    // 查看详情
    function openViewModal(btn) {
        const idea = decodeURIComponent(btn.getAttribute('data-idea'));

        document.getElementById('view_name').innerText = btn.getAttribute('data-name');
        document.getElementById('view_phone').innerText = btn.getAttribute('data-phone');
        document.getElementById('view_email').innerText = btn.getAttribute('data-email');
        document.getElementById('view_idea').innerText = idea;

        viewModal.style.display = 'block';
    }

    // 回复弹窗
    function openReplyModal(opinionId) {
        document.getElementById('replay_user').value = opinionId;
        replyModal.style.display = 'block';
    }

    // 关闭回复
    closeModal.onclick = () => replyModal.style.display = 'none';

    // 关闭查看
    closeView.onclick = () => viewModal.style.display = 'none';

    // 点击外部关闭
    window.onclick = function(event) {
        if (event.target == replyModal) {
            replyModal.style.display = 'none';
        }
        if (event.target == viewModal) {
            viewModal.style.display = 'none';
        }
    };

    // 提交回复
    replyForm.addEventListener('submit', function(e) {
        e.preventDefault();

        if (!confirm('确认提交回复？')) return;

        const formData = new FormData(replyForm);

        fetch('/replay/useropi', {
            method: 'POST',
            body: formData
        })
        .then(res => {
            if (res.ok) {
                alert('回复成功');
                replyModal.style.display = 'none';
                replyForm.reset();
                fetchOpinions('');
            } else if (res.status === 400) {
                alert("请当天回复");
            } else {
                throw new Error();
            }
        })
        .catch(() => alert('提交失败'));
    });

    // 默认值
    document.getElementById('replay_date').valueAsDate = new Date();
    document.getElementById('replay_name').value = "智慧图书馆";
});