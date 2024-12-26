$(document).ready(function () {
    let page = 1;
    const itemsPerPage = 10;
    const loadingIndicator = $('<div id="loading-indicator" class="text-center my-4"><div class="spinner-border" role="status"><span class="visually-hidden">Loading...</span></div></div>');

    function fetchData(page) {
        // Mock API call with pagination
        return new Promise((resolve) => {
            setTimeout(() => {
                const data = [];
                for (let i = 0; i < itemsPerPage; i++) {
                    data.push({
                        title: `Card ${((page - 1) * itemsPerPage) + i + 1}`,
                        text: `Content for card ${((page - 1) * itemsPerPage) + i + 1}`,
                        img: "https://via.placeholder.com/300"
                    });
                }
                console.log("Fetched data:", data);
                resolve(data);
            }, 1000);
        });
    }

    function renderMasonry(data) {
        const masonryContainer = $('#masonry');
        data.forEach(item => {
            const cardHtml = `
                <div class="col-md-12 mb-4">
                    <div class="card d-flex flex-row">
                        <img src="${item.img}" class="card-img-left" alt="...">
                        <div class="card-body">
                            <h5 class="card-title">${item.title}</h5>
                            <p class="card-text">${item.text}</p>
                        </div>
                    </div>
                </div>
            `;
            masonryContainer.append(cardHtml);
        });
    }

    function loadMore() {
        $('body').append(loadingIndicator);
        fetchData(page).then(data => {
            renderMasonry(data);
            page++;
            loadingIndicator.remove();
        });
    }

    function debounce(func, wait) {
        let timeout;
        return function (...args) {
            const later = () => {
                clearTimeout(timeout);
                func.apply(this, args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    }

    $(window).scroll(debounce(function() {
        if ($(window).scrollTop() + $(window).height() >= $(document).height() - 100) {
            loadMore();
        }
    }, 200));

    // Initial load
    loadMore();

    // Back to top button
    const backToTopBtn = $('<button id="back-to-top" class="btn btn-primary">Top</button>');
    $('body').append(backToTopBtn);
    backToTopBtn.hide();

    $(window).scroll(function() {
        if ($(window).scrollTop() > 300) {
            backToTopBtn.fadeIn();
        } else {
            backToTopBtn.fadeOut();
        }
    });

    backToTopBtn.click(function() {
        $('html, body').animate({ scrollTop: 0 }, '300');
    });

    // Mock login state
    let isLoggedIn = false;
    const mockUsername = "user";
    const mockPassword = "pass";

    function updateLoginState() {
        if (isLoggedIn) {
            $('#login-btn').hide();
            $('#user-info').removeClass('d-none');
            $('#username').text(mockUsername);
        } else {
            $('#login-btn').show();
            $('#user-info').addClass('d-none');
        }
    }

    $('#login-btn').click(function() {
        $('#loginModal').modal('show');
    });

    $('#submit-login').click(function() {
        const username = $('#login-username').val();
        const password = $('#login-password').val();
        if (username === mockUsername && password === mockPassword) {
            isLoggedIn = true;
            updateLoginState();
            $('#loginModal').modal('hide');
        } else {
            alert("Invalid credentials");
        }
    });

    $('#logout-btn').click(function() {
        isLoggedIn = false;
        updateLoginState();
    });

    $('#subscribe-btn').click(function() {
        $('#subscribeModal').modal('show');
    });

    $('#submit-subscription').click(function() {
        const url = $('#subscription-url').val();
        if (url) {
            console.log("Subscribed to:", url);
            $('#subscribeModal').modal('hide');
        } else {
            alert("Please enter a URL");
        }
    });

    // Initial login state update
    updateLoginState();
});
