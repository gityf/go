namespace go demo.rpc
namespace java demo.rpc
 
// ���Է���
service RpcService {
 
    // ����Զ�̵���
    list<string> funCall(1:i64 callTime, 2:string funCode, 3:map<string, string> paramMap),
 
}