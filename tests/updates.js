
use MoniTalks-Social 

# GMcD on Ocean
db.post.findOne({ objectId: new BinData(0, UUID("d29efbcb-10e9-4e3d-9ee3-dddab7e0fddd").base64()) })
db.post.updateOne(
                    { objectId: new BinData(0, UUID("d29efbcb-10e9-4e3d-9ee3-dddab7e0fddd").base64()) },
                    { $set: {collectiveId: new BinData(0, UUID("a7aaabc9-4053-4596-9e51-37a2295fb6a9").base64()) }}
                  )

# ROBIN ON BETA 
db.post.findOne({ objectId: new BinData(0, UUID("fdfb4daf-52e1-46b4-b05c-492b7f16e9af").base64()) })
db.post.updateOne(
                    { objectId: new BinData(0, UUID("fdfb4daf-52e1-46b4-b05c-492b7f16e9af").base64()) },
                    { $set: {collectiveId: new BinData(0, UUID("6063afd1-5024-43f0-a070-d68a54d86ec5").base64()) }}
                )

db.post.findOne({ objectId: new BinData(0, UUID("711dddc4-45be-4fdd-9ff4-c3be4dd3d789").base64()) })
db.post.updateOne(
                    { objectId: new BinData(0, UUID("711dddc4-45be-4fdd-9ff4-c3be4dd3d789").base64()) },
                    { $set: {collectiveId: new BinData(0, UUID("6063afd1-5024-43f0-a070-d68a54d86ec5").base64()) }}
                )

db.post.findOne({ objectId: new BinData(0, UUID("1e51d6a7-151c-4293-be8c-545ca30323fd").base64()) })
db.post.updateOne(
                    { objectId: new BinData(0, UUID("1e51d6a7-151c-4293-be8c-545ca30323fd").base64()) },
                    { $set: {collectiveId: new BinData(0, UUID("6063afd1-5024-43f0-a070-d68a54d86ec5").base64()) }}
                )
