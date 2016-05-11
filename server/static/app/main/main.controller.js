'use strict';

angular.module('aptWebApp')
    .controller('MainCtrl', function($scope, $translate) {
        $scope.changeLanguage = function(langKey) {
            $translate.use(langKey);
        };
    });
