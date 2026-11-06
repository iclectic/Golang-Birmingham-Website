// Progressive enhancement only: the site works fully without JavaScript.
(function () {
  'use strict';

  var menu = document.querySelector('.mobile-nav');
  if (!menu) return;

  function close() { menu.removeAttribute('open'); }

  document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape' && menu.hasAttribute('open')) {
      close();
      menu.querySelector('summary').focus();
    }
  });

  document.addEventListener('click', function (e) {
    if (menu.hasAttribute('open') && !menu.contains(e.target)) close();
  });

  menu.querySelectorAll('a').forEach(function (link) {
    link.addEventListener('click', close);
  });

  window.matchMedia('(min-width: 960px)').addEventListener('change', function (mq) {
    if (mq.matches) close();
  });
})();
