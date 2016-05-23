'use strict';

angular.module('aptWebApp')
    .directive("drawScreen", function() {
        return {
            restrict: "A",
            link: function(scope, element) {
                scope.drawElement=element;
            }
        };
    });
