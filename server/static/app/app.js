'use strict';

angular.module('aptWebApp', [
    'ngCookies',
    'ngResource',
    'ngSanitize',
    'ui.router',
    'validation.match',
    'ngFileUpload',
    'ngAnimate'
  ])
  .config(function($urlRouterProvider, $locationProvider) {
    $urlRouterProvider
      .otherwise('/table/all');

    $locationProvider.html5Mode(true);
  });
