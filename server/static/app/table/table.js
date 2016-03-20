'use strict';

angular.module('aptWebApp')
  .config(function($stateProvider) {
    $stateProvider
      .state('table', {
        url: '/',
        templateUrl: 'static/app/table/table.html',
        controller: 'TableCtrl'
      });
  });
