# Lista de tarefas

Ordem pensada para tirar a engine do estado de protótipo e colocá-la em combate jogável. Cada item abaixo bloqueia o próximo em algum nível.

1. Finalizar a aplicação de hit no runtime
   hoje a engine só detecta sobreposição; sem aplicar dano, reação e estado, não existe luta de fato.

2. Implementar o sistema de vida e reação ao golpe
   HP, hitstun, knockback, knockup e knockdown são o que transforma um acerto em consequência real.

3. Implementar bloqueio/guard
   defesa é uma mecânica básica de fighting game e evita que todo contato vire hit garantido.

4. Demais dados de framedata estarem ligados a lógica
   damage, blockstun, invincible, armor e demais campos já existem, mas ainda não governam o runtime.

5. Fechar o fluxo de round, KO e fim de partida
   sem condição de vitória e transição de match, a partida não tem começo, meio e fim.

6. Criar HUD de luta
   o jogador precisa ler vida, estado e ritmo da luta sem depender do debug.

7. Validar o ciclo completo com um personagem de teste
   depois de integrar tudo, é preciso um cenário simples para provar que o combate funciona de ponta a ponta.

Prioridade atual

1. Terminar a implementação de hits
   Isso é o bloqueio principal: sem fechamento do fluxo de hit, a luta continua só no papel.

2. Aplicar dano, hurt e framedata nos dois players
   O acerto precisa mover o estado real dos dois lados, não só acusar overlap.

3. Criar um ataque e uma animação de receber dano
   Precisamos de um ciclo mínimo para testar ataque, reação e transição de animação de ponta a ponta.

4. Consolidar o switch de animação de hurt para o player atingido
   Depois que o ciclo básico existir, a troca de estado precisa ficar confiável e previsível.

