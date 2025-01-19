$(document).ready(function () {
    let page = 1;
    const itemsPerPage = 10;
    let isLoading = false;
    const loadingIndicator = $('<div id="loading-indicator" class="text-center my-4"><div class="spinner-border" role="status"><span class="visually-hidden">Loading...</span></div></div>');
    const noMoreContentMessage = $('<div id="no-more-content" class="text-center my-4">没有更多内容</div>');

    function fetchData(page) {
        return $.ajax({
            url: '/buzz/content/list',
            method: 'POST',
            data: JSON.stringify({ PageIndex: page, PageSize: itemsPerPage }),
            headers: {
                'Authorization': 'Bearer ' + getToken(),
                'Content-Type': 'application/json'
            },
            success: function (data) {
                return data.contents;
            },
            error: function () {
                alert("Failed to fetch data");
            }
        });
    }

    function renderMasonry(contents) {
        const masonryContainer = $('#masonry');
        contents.forEach(item => {
            const cardHtml = `
                <div>
                    <div class="card d-flex flex-row" data-credit="${item.content_credit}">
                        <div class="card-img-left-wrapper">
                            <img src="${item.thumbnail}" class="card-img-left" alt="...">
                            <div class="media-length">${item.length}</div>
                        </div>
                        <div class="card-body">
                            <h5 class="card-title text-truncate">${item.name}</h5>
                            <p class="card-text text-truncate">${item.channel_name}</p>
                            <p class="card-text text-truncate">${item.published_time}</p>
                        </div>
                    </div>
                </div>
            `;
            const cardElement = $(cardHtml);
            cardElement.click(function () {
                playAudio(item);
            });
            masonryContainer.append(cardElement);
        });
    }

    function loadMore() {
        if (isLoading) return;
        isLoading = true;
        $('body').append(loadingIndicator);
        fetchData(page).then(data => {
            if (data.contents.length === 0) {
                $('body').append(noMoreContentMessage);
            } else {
                renderMasonry(data.contents);
                page++;
            }
            loadingIndicator.remove();
            isLoading = false;
        });
    }

    function fetchSub() {
        return $.ajax({
            url: '/buzz/subscription/list',
            method: 'GET',
            headers: {
                'Authorization': 'Bearer ' + getToken()
            },
            success: function (data) {
                return data.subscriptions;
            },
            error: function () {
                alert("Failed to fetch subscriptions");
            }
        });
    }

    function renderSubscriptions(data) {
        const sidebar = $('#sidebar .sidebar-content');
        sidebar.empty();
        data.subscriptions.forEach(sub => {
            const subHtml = `
                <div class="sidebar-item d-flex align-items-center">
                    <img src="${sub.channel_thumbnail}" alt="${sub.channel_name}" class="sidebar-thumbnail">
                    <p class="mb-0 ms-2">${sub.channel_name}</p>
                </div>
            `;
            sidebar.append(subHtml);
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

    $(window).scroll(debounce(function () {
        if ($(window).scrollTop() + $(window).height() >= $(document).height() - 50) {
            loadMore();
        }
    }, 100));

    // Initial load
    loadMore();
    fetchSub().then(renderSubscriptions);

    var aplayer = null;

    function playAudio(data) {
        aplayer = new APlayer({
            container: document.getElementById('aplayer'),
            audio: {
                name: data.name,
                artist: data.channel_name,
                url: `/buzz/content/stream/` + data.credit,
                cover: data.thumbnail,
                listFolded: true,
                theme: '#b7daff',
            }
        });
        aplayer.play();
    }

    // Back to top button
    const backToTopBtn = $('<button id="back-to-top" class="btn btn-primary">Top</button>');
    $('body').append(backToTopBtn);
    backToTopBtn.hide();

    $(window).scroll(function () {
        if ($(window).scrollTop() > 300) {
            backToTopBtn.fadeIn();
        } else {
            backToTopBtn.fadeOut();
        }
    });

    backToTopBtn.click(function () {
        $('html, body').animate({ scrollTop: 0 }, '300');
    });

    // Function to save JWT token
    function saveToken(token) {
        localStorage.setItem('jwtToken', token);
    }

    // Function to get JWT token
    function getToken() {
        return localStorage.getItem('jwtToken');
    }

    // Function to remove JWT token
    function removeToken() {
        localStorage.removeItem('jwtToken');
    }

    // Function to update login state
    function updateLoginState() {
        const token = getToken();
        if (token) {
            $.ajax({
                url: '/auth/current_user',
                method: 'GET',
                headers: {
                    'Authorization': 'Bearer ' + token
                },
                success: function (data) {
                    $('#login-btn').hide();
                    $('#user-info').removeClass('d-none');
                    $('#username').text(data.username);
                },
                error: function () {
                    removeToken();
                    window.location.href = '/static/login.html'; // Redirect to login page
                }
            });
        } else {
            window.location.href = '/static/login.html'; // Redirect to login page
        }
    }

    // Login button click event
    $('#login-btn').click(function () {
        $('#loginModal').modal('show');
    });

    // Submit login button click event
    $('#submit-login').click(function () {
        const username = $('#login-username').val();
        const password = $('#login-password').val();
        $.ajax({
            url: '/auth/login',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ username: username, password: password }),
            success: function (data) {
                saveToken(data.token);
                updateLoginState();
                $('#loginModal').modal('hide');
            },
            error: function () {
                alert("Invalid credentials");
            }
        });
    });

    // Logout button click event
    $('#logout-btn').click(function () {
        const token = getToken();
        if (token) {
            $.ajax({
                url: '/auth/logout',
                method: 'POST',
                headers: {
                    'Authorization': 'Bearer ' + token
                },
                success: function () {
                    removeToken();
                    updateLoginState();
                },
                error: function () {
                    alert("Failed to logout");
                }
            });
        }
    });

    // Initial login state update
    updateLoginState();

    $('#subscribe-btn').click(function () {
        $('#subscribeModal').modal('show');
    });

    $('#submit-subscription').click(function () {
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

    // Ensure the container and search input have the sidebar-collapsed class initially
    $('.container').addClass('sidebar-collapsed');
    $('.search').addClass('sidebar-collapsed');

    // Sidebar toggle functionality
    $('#sidebar-toggle').click(function () {
        $('#sidebar').toggleClass('collapsed');
        $('#sidebar').toggleClass('sidebar-expanded');
        $('.container').toggleClass('sidebar-expanded sidebar-collapsed');
        $('.search').toggleClass('sidebar-expanded sidebar-collapsed');

        // Change the icon based on the sidebar state
        if ($('#sidebar').hasClass('collapsed')) {
            $('#sidebar-toggle i').removeClass('fa-chevron-left').addClass('fa-bars');
        } else {
            $('#sidebar-toggle i').removeClass('fa-bars').addClass('fa-chevron-left');
        }
        adjustMasonryItemWidth();
    });

    function adjustMasonryItemWidth() {
        const masonryContainer = $('#masonry');
        if ($('#sidebar').hasClass('collapsed')) {
            masonryContainer.find('.card').css('width', 'calc(100% - 60px)');
        } else {
            masonryContainer.find('.card').css('width', 'calc(100% - 250px)');
        }
    }

    // Initial load
    loadMore();
    fetchSub().then(renderSubscriptions);
    adjustMasonryItemWidth();
});
