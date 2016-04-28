'use strict';

angular.module('aptWebApp')
    .controller('DeviceDetailCtrl', function($scope, $http, $stateParams) {
        $scope.deviceID = $stateParams.id;
        $scope.nodeIP = $stateParams.ip;
        $scope.device={};
        $scope.refresh = function() {
            $http.get('device/' + $scope.nodeIP + "/" + $scope.deviceID).then(response => {
                if (response) {
                    $scope.device = response.data;
                }
            });
        }

        $scope.refresh();
    });
