#include <iostream>
#include <cstdio>
#include <cstring>
using namespace std;
int a[10];
int main(){
	int n;
	scanf("%d",&n);
	for(int i = 0; i < 1<<n; i++){
		memset(a,0,sizeof(a));
		for(int j = 0; j < n; j++){
			if(i & (1<<j)){
				a[j] = 1;
			}
		}
		for(int j = 0 ; j < n; j++){
			if(a[j]){
				cout<<"./client subjob sub_1.json"<<endl;
			}else{
				cout<<"./client subjob sub_2.json"<<endl;
			}
		}
		cout<<endl;		
	}
}
