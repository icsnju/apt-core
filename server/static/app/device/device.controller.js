'use strict';

angular.module('aptWebApp')
    .controller('DeviceDetailCtrl', function($scope, $http, $stateParams) {
        $scope.deviceID = $stateParams.id;
        $scope.nodeIP = '';
        $scope.device = {};
        //the status of device
        $scope.getState = function(state) {
            if (state == 0) {
                return 'busy';
            } else {
                return 'free';
            }
        };

        //send device button events to websocket
        $scope.deviceButton = function(kind) {
            if ($scope.screenWS) {
                $scope.screenWS.send(kind)
            }
        }

        $http.get('device/ip/' + $scope.deviceID).then(function(response) {
            if (response) {
                $scope.nodeIP = response.data;
            }
        }, function(response) {

        });

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
