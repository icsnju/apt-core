'use strict';

angular.module('aptWebApp')
  .config(function($stateProvider, $urlRouterProvider) {
    $stateProvider
      .state('submit', {
        url: '/submit',
        templateUrl: 'static/app/submit/submit.html',
        controller: 'SubmitCtrl'
      })
      .state('submit.frame', {
        url: '/frame',
        templateUrl: 'static/app/submit/submit-frame.html'
      })
      .state('submit.selector', {
        url: '/selector',
        templateUrl: 'static/app/submit/submit-selector.html'
      })
      .state('submit.ok', {
        url: '/ok',
        templateUrl: 'static/app/submit/submit-ok.html'
      });
  });
