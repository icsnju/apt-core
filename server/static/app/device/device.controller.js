'use strict';

angular.module('aptWebApp')
    .controller('DeviceDetailCtrl', function($scope, $http, $stateParams) {
        $scope.deviceID = $stateParams.id;
        $scope.device = {};
        $scope.refresh = function() {
            $http.get('device/' + $scope.deviceID).then(function(response) {
                if (response) {
                    $scope.device = response.data;
                }
            }, function(response) {
                //console.log(response)
            });
        }

        $scope.refresh();
    });
