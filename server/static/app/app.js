'use strict';

angular.module('aptWebApp', [
        'ngCookies',
        'ngResource',
        'ngSanitize',
        'ui.router',
        'validation.match',
        'ngFileUpload',
        'ngAnimate',
        'pascalprecht.translate',
        'luegg.directives'
    ])
    .config(function($urlRouterProvider, $locationProvider) {
        $urlRouterProvider
            .otherwise('/jobs');

        $locationProvider.html5Mode(true);
    });
