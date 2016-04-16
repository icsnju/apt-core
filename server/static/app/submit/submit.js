'use strict';

angular.module('aptWebApp')
  .config(function($stateProvider, $urlRouterProvider) {
    $stateProvider
      .state('main.submit', {
        url: '/submit',
        templateUrl: 'static/app/submit/submit.html',
        controller: 'SubmitCtrl'
      })
      .state('main.submit.frame', {
        url: '/frame',
        templateUrl: 'static/app/submit/submit.frame.html'
      })
      .state('main.submit.selector', {
        url: '/selector',
        templateUrl: 'static/app/submit/submit.selector.html'
      })
      .state('main.submit.ok', {
        url: '/ok',
        templateUrl: 'static/app/submit/submit.ok.html'
      });
  });
